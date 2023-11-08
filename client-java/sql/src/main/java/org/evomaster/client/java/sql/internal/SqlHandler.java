package org.evomaster.client.java.sql.internal;

import net.sf.jsqlparser.expression.Expression;
import net.sf.jsqlparser.expression.ExpressionVisitor;
import net.sf.jsqlparser.expression.ExpressionVisitorAdapter;
import net.sf.jsqlparser.parser.CCJSqlParserUtil;
import net.sf.jsqlparser.schema.Column;
import net.sf.jsqlparser.statement.Statement;
import org.evomaster.client.java.controller.api.dto.database.execution.ExecutionDto;
import org.evomaster.client.java.controller.api.dto.database.execution.SqlExecutionLogDto;
import org.evomaster.client.java.controller.api.dto.database.schema.DbSchemaDto;
import org.evomaster.client.java.sql.distance.advanced.AdvancedDistance;
import org.evomaster.client.java.sql.distance.advanced.driver.SqlDriver;
import org.evomaster.client.java.sql.distance.advanced.driver.cache.ConcurrentCache;
import org.evomaster.client.java.sql.distance.standard.StandardDistance;
import org.evomaster.client.java.utils.SimpleLogger;

import java.sql.Connection;
import java.util.*;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.CopyOnWriteArrayList;

import static java.util.Objects.isNull;
import static org.evomaster.client.java.sql.distance.advanced.driver.SqlDriver.createSqlDriver;
import static org.evomaster.client.java.sql.internal.ParserUtils.*;

/**
 * Class used to act upon SQL commands executed by the SUT
 */
public class SqlHandler {

    private final static Set<String> booleanConstantNames = Collections.unmodifiableSet(
            new LinkedHashSet<>(Arrays.asList("t", "true", "f", "false", "yes", "y", "no", "n", "on", "off", "unknown"))
    );

    private final TaintHandler taintHandler;

    /**
     * Computing heuristics on SQL is expensive, as we need to run
     * further queries. So, we buffer them, and execute them only
     * if needed (ie, lazy initialization)
     */
    private final List<String> buffer;

    /**
     * The heuristics based on the SQL execution
     */
    private final List<PairCommandDistance> distances;

    //see ExecutionDto
    private final Map<String, Set<String>> queriedData;
    private final Map<String, Set<String>> updatedData;
    private final Map<String, Set<String>> insertedData;
    private final Map<String, Set<String>> failedWhere;
    private final List<String> deletedData;
    private final List<SqlExecutionLogDto> executedInfo;

    private int numberOfSqlCommands;

    private volatile Connection connection;

    private volatile boolean calculateHeuristics;

    private volatile boolean extractSqlExecution;

    private volatile boolean advancedHeuristics;

    private volatile SqlDriver sqlDriver;

    /**
     * WARNING: in general we shouldn't use mutable DTO as internal data structures.
     * But, here, what we need is very simple (just checking for names).
     */
    private volatile DbSchemaDto schema;

    public SqlHandler(TaintHandler taintHandler) {

        this.taintHandler = taintHandler;

        buffer = new CopyOnWriteArrayList<>();
        distances = new ArrayList<>();
        queriedData = new ConcurrentHashMap<>();
        updatedData = new ConcurrentHashMap<>();
        insertedData = new ConcurrentHashMap<>();
        failedWhere = new ConcurrentHashMap<>();
        deletedData = new CopyOnWriteArrayList<>();
        executedInfo = new CopyOnWriteArrayList<>();

        calculateHeuristics = true;
        numberOfSqlCommands = 0;
    }

    public void reset() {
        buffer.clear();
        distances.clear();
        queriedData.clear();
        updatedData.clear();
        insertedData.clear();
        failedWhere.clear();
        deletedData.clear();
        executedInfo.clear();

        numberOfSqlCommands = 0;

        sqlDriver = null;
    }

    public void setConnection(Connection connection) {
        this.connection = connection;
    }

    public void setSchema(DbSchemaDto schema) {
        this.schema = schema;
    }

    /**
     * handle executed sql info
     *
     * @param sql to be handled
     */
    public void handle(SqlExecutionLogDto sql) {
        executedInfo.add(sql);
        handle(sql.command);
    }

    public void handle(String sql) {
        Objects.requireNonNull(sql);

        if (!calculateHeuristics && !extractSqlExecution) {
            return;
        }

        numberOfSqlCommands++;

        if (!ParserUtils.canParseSqlStatement(sql)) {
            SimpleLogger.warn("Cannot parse SQL statement: " + sql);
            return;
        }

        buffer.add(sql);

        if (isSelect(sql)) {
            mergeNewData(queriedData, ColumnTableAnalyzer.getSelectReadDataFields(sql));
        } else if (isDelete(sql)) {
            deletedData.addAll(ColumnTableAnalyzer.getDeletedTables(sql));
        } else if (isInsert(sql)) {
            mergeNewData(insertedData, ColumnTableAnalyzer.getInsertedDataFields(sql));
        } else if (isUpdate(sql)) {
            mergeNewData(updatedData, ColumnTableAnalyzer.getUpdatedDataFields(sql));
        }

    }

    public ExecutionDto getExecutionDto() {

        if (!calculateHeuristics && !extractSqlExecution) {
            return null;
        }

        ExecutionDto executionDto = new ExecutionDto();
        executionDto.queriedData.putAll(queriedData);
        executionDto.failedWhere.putAll(failedWhere);
        executionDto.insertedData.putAll(insertedData);
        executionDto.updatedData.putAll(updatedData);
        executionDto.deletedData.addAll(deletedData);
        executionDto.numberOfSqlCommands = this.numberOfSqlCommands;
        executionDto.sqlExecutionLogDtoList.addAll(executedInfo);
        return executionDto;
    }

    /**
     * compute (SELECT, DELETE and UPDATE) sql distance for sql commands which exists in [buffer]
     * Note that we skip `SELECT 1` (typically for testing sql connection) since its distance is 0
     *
     * @return a list of heuristics for sql commands
     */
    public List<PairCommandDistance> getDistances() {

        if (connection == null || !calculateHeuristics) {
            return distances;
        }


        buffer.forEach(sql -> {
                    if (!isSelectOne(sql) && (isSelect(sql) || isDelete(sql) || isUpdate(sql))) {
                        double dist;
                        try {
                            dist = computeDistance(sql);
                        } catch (Exception e) {
                            SimpleLogger.error("FAILED TO COMPUTE HEURISTICS FOR SQL: " + sql);
                            dist = Double.MAX_VALUE;
                            //assert false; //TODO put back once we update JSqlParser
                            //return;
                        }
                        distances.add(new PairCommandDistance(sql, dist));
                    }
                });

        //side effects on buffer is not important, as it is just a cache
        buffer.clear();

        return distances;
    }


    private Double computeDistance(String command) {

        if (connection == null) {
            throw new IllegalStateException("Trying to calculate SQL distance with no DB connection");
        }

        Statement statement;

        try {
            statement = CCJSqlParserUtil.parse(command);
        } catch (Exception e) {
            SimpleLogger.uniqueWarn("Cannot handle command: " + command + "\n" + e);
            return Double.MAX_VALUE;
        }


        Map<String, Set<String>> columns = extractColumnsInvolvedInWhere(statement);

        /*
            even if columns.isEmpty(), we need to check if any data was present
         */

        double dist;
        if (columns.isEmpty()) {
            //TODO check if table(s) not empty, and give >0 otherwise
            dist = 0;
        } else {
            if(isAdvancedHeuristics()){
                if(isNull(sqlDriver)){
                    sqlDriver = createSqlDriver(connection, new ConcurrentCache());
                }
                AdvancedDistance advancedDistance = new AdvancedDistance(sqlDriver);
                dist = advancedDistance.calculate(command);
            } else {
                StandardDistance standardDistance = new StandardDistance(connection, schema, taintHandler);
                dist = standardDistance.calculateDistance(command, columns);
            }
        }

        if (dist > 0) {
            mergeNewData(failedWhere, columns);
        }

        return dist;
    }

    /**
     * Check the fields involved in the WHERE clause (if any).
     * Return a map from table name to column names of the involved fields.
     */
    public Map<String, Set<String>> extractColumnsInvolvedInWhere(Statement statement) {

        /*
           TODO
           following does not handle the case of sub-selects involving other
           tables... but likely that is not something we need to support right now
         */

        Map<String, Set<String>> data = new HashMap<>();

        // move getWhere before SqlNameContext, otherwise null where would cause exception in new SqlNameContext
        Expression where = ParserUtils.getWhere(statement);
        if (where == null) {
            return data;
        }

        SqlNameContext context = new SqlNameContext(statement);
        if (schema != null) {
            context.setSchema(schema);
        }

        ExpressionVisitor visitor = new ExpressionVisitorAdapter() {
            @Override
            public void visit(Column column) {

                String tn = context.getTableName(column);

                if (tn.equalsIgnoreCase(SqlNameContext.UNNAMED_TABLE)) {
                    // TODO handle it properly when ll have support for sub-selects
                    return;
                }

                String cn = column.getColumnName().toLowerCase();

                if (!context.hasColumn(tn, cn)) {

                    /*
                        This is an issue with the JsqlParser library. Until we upgrade it, or fix it if not fixed yet,
                        we use this workaround.

                        The problem is that some SQL databases do not have support for boolean types, so parser can
                        interpret constants like TRUE as column names.
                        And all databases have differences on how booleans are treated, eg.
                        - H2: TRUE, FALSE, and UNKNOWN (NULL).
                          http://www.h2database.com/html/datatypes.html#boolean_type
                        - Postgres:  true, yes, on, 1, false, no, off, 0 (as well as abbreviations like t and f)
                          https://www.postgresql.org/docs/9.5/datatype-boolean.html

                     */

                    if (booleanConstantNames.contains(cn)) {
                        //case in which a boolean constant is wrongly treated as a column name.
                        //TODO not sure what we can really do here without modifying the parser
                    } else {
                        SimpleLogger.warn("Cannot find column '" + cn + "' in table '" + tn + "'");
                    }
                    return;
                }

                data.putIfAbsent(tn, new HashSet<>());
                Set<String> set = data.get(tn);
                set.add(cn);
            }
        };

        where.accept(visitor);

        return data;
    }

    private static void mergeNewData(
            Map<String, Set<String>> current,
            Map<String, Set<String>> toAdd
    ) {

        for (Map.Entry<String, Set<String>> e : toAdd.entrySet()) {
            String key = e.getKey();
            Set<String> values = e.getValue();

            Set<String> existing = current.get(key);

            if (existing != null && existing.contains("*")) {
                //nothing to do
                continue;
            }

            if (existing == null) {
                existing = new HashSet<>(values);
                current.put(key, existing);
            } else {
                existing.addAll(values);
            }

            if (existing.size() > 1 && existing.contains("*")) {
                /*
                    remove unnecessary columns, as anyway we take
                    everything with *
                 */
                existing.clear();
                existing.add("*");
            }
        }
    }

    public boolean isCalculateHeuristics() {
        return calculateHeuristics;
    }

    public boolean isExtractSqlExecution() {
        return extractSqlExecution;
    }

    public void setCalculateHeuristics(boolean calculateHeuristics) {
        this.calculateHeuristics = calculateHeuristics;
    }

    public void setExtractSqlExecution(boolean extractSqlExecution) {
        this.extractSqlExecution = extractSqlExecution;
    }

    public boolean isAdvancedHeuristics() {
        return advancedHeuristics;
    }

    public void setAdvancedHeuristics(boolean advancedHeuristics) {
        this.advancedHeuristics = advancedHeuristics;
    }
}