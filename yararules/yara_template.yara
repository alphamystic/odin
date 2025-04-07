/*
  This is a yara template for autogeneration of yara rules via odin

*/


rule Template

{

  meta:
    description = "Rule Description"
    author = "Author name"
    source = "SOurce or report if any"
    hash = "Not sure if hash of file or rule"

  strings:
    $str1 = "somestring" no case
    $str2 = "otherstring" fullword ascii
    $str3 = "any string" ascii wide

  condition:
    all of them

}
