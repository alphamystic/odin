rule url_and_ip_windows_binaries
{
  meta:
    description = "Optimized rule for Windows binaries (.exe, .dll, .msi, .msc) with URLs and IPs"
    author = "Samuel Odhiambo"

  strings:
    $url_http = "http" wide ascii fullword
    $url_https = "https" wide ascii fullword
    $ip_pattern = /([0-9]{1,3}\.){3}[0-9]{1,3}/ ascii

  condition:
    uint16(0) == 0x5a4d and
    (#url_http + #url_https + #ip_pattern) >= 3
}


rule url_and_ip_linux_binaries
{
  meta:
    description = "Optimized rule for Linux binaries (.so and ELF binaries) with URLs and IPs"
    author = "Samuel Odhiambo"

  strings:
    $url_http = "http" wide ascii
    $url_https = "https" wide ascii
    $ip_pattern = /([0-9]{1,3}\.){3}[0-9]{1,3}/ ascii

  condition:
    uint32(0) == 0x7f454c46 and
    (#url_http + #url_https + #ip_pattern) >= 3
}
