package org.evomaster.core.output


class TestSuiteFileName(
        /**
         * This can be in the form
         * foo.bar.X
         */
        val name: String) {


    fun getPackage() : String{
        if(! name.contains('.')){
            return ""
        }

        return name.substring(0, name.lastIndexOf('.'))
    }

    fun hasPackage() = ! getPackage().isBlank()


    fun getClassName(format: OutputFormat): String{
        if(! hasPackage()){
            when {
                format.isGo() -> return name.split('_').joinToString("", transform = { it.replaceFirstChar { it.uppercase() }} )
            }

            return name
        }

        return name.substring(name.lastIndexOf('.') + 1, name.length)
    }


    fun getAsPath(format: OutputFormat) : String{

        //TODO what about C#? is it a behavior we want there as well
        var base = if(format.isJavaOrKotlin()) name.replace('.', '/') else name
        when {
            format.isGo() -> base = base.lowercase()
        }

        return base + when{
            format.isJava() -> ".java"
            format.isKotlin() -> ".kt"
            format.isJavaScript() -> ".js"
            format.isCsharp() -> ".cs"
            format.isGo() -> ".go"
            else -> throw IllegalStateException("Unsupported format $format")
        }
    }
}