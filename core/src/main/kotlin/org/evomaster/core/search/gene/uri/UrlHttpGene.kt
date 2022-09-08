package org.evomaster.core.search.gene.uri

import org.evomaster.core.output.OutputFormat
import org.evomaster.core.search.gene.*
import org.evomaster.core.search.gene.network.InetGene
import org.evomaster.core.search.service.AdaptiveParameterControl
import org.evomaster.core.search.service.Randomness
import org.evomaster.core.search.service.mutator.genemutation.AdditionalGeneMutationInfo
import org.evomaster.core.search.service.mutator.genemutation.SubsetGeneSelectionStrategy
import java.net.URL


/**
 * Form:
 * <user>:<password>@<host>:<port>/<url-path>
 */
class UrlHttpGene(
    name: String,
    val scheme : EnumGene<String> = EnumGene("scheme", listOf("http","https")),
    //TODO authority: <user>:<password>@
    val host: ChoiceGene<Gene> = ChoiceGene("host", listOf(
        InetGene("ipv4"),
        StringGene("hostname", invalidChars = invalidChars)
    )),
    val port : OptionalGene = OptionalGene("port", IntegerGene("port", min=0,max=65535)),
    val path : ArrayGene<StringGene> = ArrayGene("path",
            minSize = 1,
            maxSize = 3,
            openingTag = "",
            closingTag = "",
            separatorTag = "/",
            template = StringGene("path", invalidChars = invalidChars)
        )
    //TODO query params ?x=y and fragment #
) : CompositeFixedGene(name, mutableListOf(scheme,host,port,path)) {

    companion object{
         val invalidChars = listOf('*','+','\\','/','#','$','!','?','[',']','{','}','(',')','\'','"')
    }


    override fun copyContent(): Gene {
        return UrlHttpGene(
            name,
            scheme.copy() as EnumGene<String>,
            host.copy() as ChoiceGene<Gene>,
            port.copy() as OptionalGene,
            path.copy() as ArrayGene<StringGene>
        )
    }

    override fun isLocallyValid(): Boolean {
        return getViewOfChildren().all { it.isLocallyValid() }
                && try{URL(getValueAsRawString()); true}catch (e: Exception){false}
    }

    override fun randomize(randomness: Randomness, tryToForceNewValue: Boolean) {
        getViewOfChildren().forEach { it.randomize(randomness, tryToForceNewValue) }
    }

    override fun candidatesInternalGenes(
        randomness: Randomness,
        apc: AdaptiveParameterControl,
        selectionStrategy: SubsetGeneSelectionStrategy,
        enableAdaptiveGeneMutation: Boolean,
        additionalGeneMutationInfo: AdditionalGeneMutationInfo?
    ): List<Gene> {

        return innerGene().filter { it.isMutable() }
    }

    override fun innerGene(): List<Gene> {
        return listOf(scheme, host, port, path)
    }

    override fun getValueAsPrintableString(
        previousGenes: List<Gene>,
        mode: GeneUtils.EscapeMode?,
        targetFormat: OutputFormat?,
        extraCheck: Boolean
    ): String {

        val s = scheme.getValueAsRawString()
        val h = host.getValueAsRawString()
        val p = if(port.isActive) ":${port.gene.getValueAsRawString()}" else ""
        val e = path.getValueAsRawString().replace("\"","")

        return "$s://$h$p/$e"
    }

    override fun copyValueFrom(other: Gene) {
        if (other !is UrlHttpGene) {
            throw IllegalArgumentException("Invalid gene type ${other.javaClass}")
        }
        scheme.copyValueFrom(other.scheme)
        host.copyValueFrom(other.host)
        port.copyValueFrom(other.port)
        path.copyValueFrom(other.path)
    }

    override fun containsSameValueAs(other: Gene): Boolean {
        if (other !is UrlHttpGene) {
            throw IllegalArgumentException("Invalid gene type ${other.javaClass}")
        }
        return scheme.containsSameValueAs(other.scheme)
                && host.containsSameValueAs(other.host)
                && port.containsSameValueAs(other.port)
                && path.containsSameValueAs(other.path)
    }

    override fun bindValueBasedOn(gene: Gene): Boolean {
        return false
    }


}