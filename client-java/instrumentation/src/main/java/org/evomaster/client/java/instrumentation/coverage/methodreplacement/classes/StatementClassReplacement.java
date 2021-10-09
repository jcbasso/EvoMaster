package org.evomaster.client.java.instrumentation.coverage.methodreplacement.classes;

import org.evomaster.client.java.instrumentation.SqlInfo;
import org.evomaster.client.java.instrumentation.coverage.methodreplacement.MethodReplacementClass;
import org.evomaster.client.java.instrumentation.coverage.methodreplacement.Replacement;
import org.evomaster.client.java.instrumentation.shared.ReplacementType;
import org.evomaster.client.java.instrumentation.staticstate.ExecutionTracer;
import org.evomaster.client.java.utils.SimpleLogger;

import java.sql.ResultSet;
import java.sql.SQLException;
import java.sql.Statement;

public class StatementClassReplacement implements MethodReplacementClass {

    @Override
    public Class<?> getTargetClass() {
        return Statement.class;
    }

    private static void handleSql(String sql){
        /*
            TODO need to provide proper info data here.
            Bit tricky, need to check actual DB implementations, see:
            https://stackoverflow.com/questions/867194/java-resultset-how-to-check-if-there-are-any-results/15750832#15750832

            Anyway, not needed till we support constraint solving for DB data, as then
            we can skip the branch distance computation
         */
        SqlInfo info = new SqlInfo(sql, false, false);
        ExecutionTracer.addSqlInfo(info);
    }

    @Replacement(type = ReplacementType.TRACKER)
    public static ResultSet executeQuery(Statement caller, String sql) throws SQLException{
        handleSql(sql);
        return caller.executeQuery(sql);
    }

    @Replacement(type = ReplacementType.TRACKER)
    public static int executeUpdate(Statement caller, String sql) throws SQLException{
        handleSql(sql);
        return caller.executeUpdate(sql);
    }

    @Replacement(type = ReplacementType.TRACKER)
    public static boolean execute(Statement caller,String sql) throws SQLException{
        handleSql(sql);
        return caller.execute(sql);
    }

    @Replacement(type = ReplacementType.TRACKER)
    public static int executeUpdate(Statement caller, String sql, int autoGeneratedKeys) throws SQLException{
        handleSql(sql);
        return caller.executeUpdate(sql, autoGeneratedKeys);
    }

    @Replacement(type = ReplacementType.TRACKER)
    public static int executeUpdate(Statement caller, String sql, int columnIndexes[]) throws SQLException{
        handleSql(sql);
        return caller.executeUpdate(sql, columnIndexes);
    }


    @Replacement(type = ReplacementType.TRACKER)
    public static int executeUpdate(Statement caller, String sql, String columnNames[]) throws SQLException{
        handleSql(sql);
        return caller.executeUpdate(sql, columnNames);
    }

    @Replacement(type = ReplacementType.TRACKER)
    public static boolean execute(Statement caller, String sql, int autoGeneratedKeys) throws SQLException{
        handleSql(sql);
        return caller.execute(sql, autoGeneratedKeys);
    }

    @Replacement(type = ReplacementType.TRACKER)
    public static boolean execute(Statement caller, String sql, int columnIndexes[]) throws SQLException{
        handleSql(sql);
        return caller.execute(sql, columnIndexes);
    }

    @Replacement(type = ReplacementType.TRACKER)
    public static boolean execute(Statement caller, String sql, String columnNames[]) throws SQLException{
        handleSql(sql);
        return caller.execute(sql, columnNames);
    }

    @Replacement(type = ReplacementType.TRACKER)
    public static long executeLargeUpdate(Statement caller, String sql) throws SQLException {
        handleSql(sql);
        return caller.executeLargeUpdate(sql);
    }

    @Replacement(type = ReplacementType.TRACKER)
    public static long executeLargeUpdate(Statement caller, String sql, int autoGeneratedKeys) throws SQLException {
        handleSql(sql);
        return caller.executeLargeUpdate(sql, autoGeneratedKeys);
    }

    @Replacement(type = ReplacementType.TRACKER)
    public static long executeLargeUpdate(Statement caller, String sql, int columnIndexes[]) throws SQLException {
        handleSql(sql);
        return caller.executeLargeUpdate(sql, columnIndexes);
    }

    @Replacement(type = ReplacementType.TRACKER)
    public static long executeLargeUpdate(Statement caller, String sql, String columnNames[]) throws SQLException {
        handleSql(sql);
        return caller.executeLargeUpdate(sql, columnNames);
    }

}