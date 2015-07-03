# This module contains the contants used for Telnet
# Based on code by Jon A. Lambert,  under the Zlib license.
 
module Telnet

module Codes

  IAC = 255  # Command  - RFC 854, 855, 1123, 1143

  # 2 byte commands
  WILL = 251  # Will do option
  WONT = 252  # Wont do option
  DO   = 253  # Do option
  DONT = 254  # Dont do option


  SB = 250  # Subnegotiation begin # IAC SB <option> <parameters> IAC SE
  SE = 240  # Subnegotiation end

  # 1 byte commands
  GA    = 249  # Go Ahead
  NOP   = 241  # No-op
  BRK   = 243  # Break

  # In RFC 854
  AYT   = 246  # Are you there?
  AO    = 245  # abort output
  IP    = 244  # interrupt
  EL    = 248  # erase current line
  EC    = 247  # erase current character

  DM    = 242  # data mark - sent to demarcate end of urgent commands

  EOR   = 239 # end of record (transparent mode)
  ABORT = 238 # Abort process
  SUSP  = 237 # Suspend process
  EOF   = 236 # End of file

  # Options
  BINARY         =   0 # Transmit Binary - RFC 856
  ECHO           =   1 # Echo - RFC 857
  RCP            =   2 # Reconnection
  SGA            =   3 # Suppress Go Ahead - RFC 858
  NAMS           =   4 # Approx Message Size Negotiation
  STATUS         =   5 # Status - RFC 859
  TM             =   6 # Timing Mark - RFC 860
  RCTE           =   7 # Remote Controlled Trans and Echo - RFC 563, 726
  NAOL           =   8 # Output Line Width
  NAOP           =   9 # Output Page Size
  NAOCRD         =  10 # Output Carriage-Return Disposition - RFC 652
  NAOHTS         =  11 # Output Horizontal Tab Stops - RFC 653
  NAOHTD         =  12 # Output Horizontal Tab Disposition - RFC 654
  NAOFFD         =  13 # Output Formfeed Disposition - RFC 655
  NAOVTS         =  14 # Output Vertical Tabstops - RFC 656
  NAOVTD         =  15 # Output Vertical Tab Disposition - RFC 657
  NAOLFD         =  16 # Output Linefeed Disposition - RFC 658
  XASCII         =  17 # Extended ASCII - RFC 698
  LOGOUT         =  18 # Logout - RFC 727
  BM             =  19 # Byte Macro - RFC 735
  DET            =  20 # Data Entry Terminal - RFC 732, 1043
  SUPDUP         =  21 # SUPDUP - RFC 734, 736
  SUPDUPOUTPUT   =  22 # SUPDUP Output - RFC 749
  SNDLOC         =  23 # Send Location - RFC 779
  TTYPE          =  24 # Terminal Type - RFC 1091
  EOREC          =  25 # End of Record - RFC 885
  TUID           =  26 # TACACS User Identification - RFC 927
  OUTMRK         =  27 # Output Marking - RFC 933
  TTYLOC         =  28 # Terminal Location Number - RFC 946
  REGIME3270     =  29 # Telnet 3270 Regime - RFC 1041
  X3PAD          =  30 # X.3 PAD - RFC 1053
  NAWS           =  31 # Negotiate About Window Size - RFC 1073
  TSPEED         =  32 # Terminal Speed - RFC 1079
  LFLOW          =  33 # Remote Flow Control - RFC 1372
  LINEMODE       =  34 # Linemode - RFC 1184
  XDISPLOC       =  35 # X Display Location - RFC 1096
  ENVIRON        =  36 # Environment Option - RFC 1408
  AUTHENTICATION =  37 # Authentication Option - RFC 1416, 2941, 2942, 2943, 2951
  ENCRYPT        =  38 # Encryption Option - RFC 2946
  NEW_ENVIRON    =  39 # New Environment Option - RFC 1572
  TN3270         =  40 # TN3270 Terminal Entry - RFC 2355
  XAUTH          =  41 # XAUTH
  CHARSET        =  42 # Charset option - RFC 2066
  RSP            =  43 # Remote Serial Port
  CPCO           =  44 # COM port Control Option - RFC 2217
  SUPLECHO       =  45 # Suppress Local Echo
  TLS            =  46 # Telnet Start TLS
  KERMIT         =  47 # Kermit tranfer Option - RFC 2840
  SENDURL        =  48 # Send URL
  FORWARDX       =  49 # Forward X
  PLOGON         = 138 # Telnet Pragma Logon
  SSPI           = 139 # Telnet SSPI Logon
  PHEARTBEAT     = 140 # Telnat Pragma Heartbeat
  EXOPL          = 255 # Extended-Options-List - RFC 861

  MSDP     = 69  # Mud Server Data Protocol
  MSSP     = 70  # MUD Server Status Protocol
  COMPRESS =  85 # MCCP 1 support (broken deprecated)
  COMPRESS2 = 86 # MCCP 2 support
  MSP  = 90 # MSP  support
  MXP  = 91 # MUD eXtension Protocol (MXP)
  MSP2 = 92 # MSP2 support
    MUSIC = 0
    SOUND = 1

  ZMP = 93 # ZMP support
  AARD = 102 # Aardwolf client protocol (deprecated?)
  MULTIPLEX = 112 # Crystal client telnet multiplex 
  ATCP = 200 # Achaea telnet client protocol
  GMCP = 201 # GMCP/ATCP2 client protocol
end
