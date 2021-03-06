# This module contains the contants used for Telnet
# Based on code by Jon A. Lambert,  under the Zlib license.
 
class Telnet
  module Codes
    # Telnet commands
    TELNET_IAC = 255
    TELNET_DONT = 254
    TELNET_DO = 253
    TELNET_WONT = 252
    TELNET_WILL = 251
    TELNET_SB = 250
    TELNET_GA = 249
    TELNET_EL = 248
    TELNET_EC = 247
    TELNET_AYT = 246
    TELNET_AO = 245
    TELNET_IP = 244
    TELNET_BREAK = 243
    TELNET_DM = 242
    TELNET_NOP = 241
    TELNET_SE = 240
    TELNET_EOR = 239
    TELNET_ABORT = 238
    TELNET_SUSP = 237
    TELNET_EOF = 236

    # Telnet options.
    TELNET_TELOPT_BINARY = 0
    TELNET_TELOPT_ECHO = 1
    TELNET_TELOPT_RCP = 2
    TELNET_TELOPT_SGA = 3
    TELNET_TELOPT_NAMS = 4
    TELNET_TELOPT_STATUS = 5
    TELNET_TELOPT_TM = 6
    TELNET_TELOPT_RCTE = 7
    TELNET_TELOPT_NAOL = 8
    TELNET_TELOPT_NAOP = 9
    TELNET_TELOPT_NAOCRD = 10
    TELNET_TELOPT_NAOHTS = 11
    TELNET_TELOPT_NAOHTD = 12
    TELNET_TELOPT_NAOFFD = 13
    TELNET_TELOPT_NAOVTS = 14
    TELNET_TELOPT_NAOVTD = 15
    TELNET_TELOPT_NAOLFD = 16
    TELNET_TELOPT_XASCII = 17
    TELNET_TELOPT_LOGOUT = 18
    TELNET_TELOPT_BM = 19
    TELNET_TELOPT_DET = 20
    TELNET_TELOPT_SUPDUP = 21
    TELNET_TELOPT_SUPDUPOUTPUT = 22
    TELNET_TELOPT_SNDLOC = 23
    TELNET_TELOPT_TTYPE = 24
    TELNET_TELOPT_EOR = 25
    TELNET_TELOPT_TUID = 26
    TELNET_TELOPT_OUTMRK = 27
    TELNET_TELOPT_TTYLOC = 28
    TELNET_TELOPT_3270REGIME = 29
    TELNET_TELOPT_X3PAD = 30
    TELNET_TELOPT_NAWS = 31
    TELNET_TELOPT_TSPEED = 32
    TELNET_TELOPT_LFLOW = 33
    TELNET_TELOPT_LINEMODE = 34
    TELNET_TELOPT_XDISPLOC = 35
    TELNET_TELOPT_ENVIRON = 36
    TELNET_TELOPT_AUTHENTICATION = 37
    TELNET_TELOPT_ENCRYPT = 38
    TELNET_TELOPT_NEW_ENVIRON = 39
    TELNET_TELOPT_MSDP = 69
    TELNET_TELOPT_MSSP = 70
    TELNET_TELOPT_COMPRESS = 85
    TELNET_TELOPT_COMPRESS2 = 86
    TELNET_TELOPT_MSP = 90
    TELNET_TELOPT_MXP = 91
    TELNET_TELOPT_MSP2 = 92
    TELNET_TELOPT_MSP2_MUSIC = 0
    TELNET_TELOPT_MSP2_SOUND = 1



    TELNET_TELOPT_ZMP = 93
    TELNET_TELOPT_EXOPL = 255

    TELNET_TELOPT_MCCP2 = 86

    # TERMINAL-TYPE codes. 
    TELNET_TTYPE_IS = 0
    TELNET_TTYPE_SEND = 1
    
    # MTTS standard codes
    TELNET_MTTS_ANSI                = 1
    TELNET_MTTS_VT100               = 2
    TELNET_MTTS_UTF8                = 4
    TELNET_MTTS_256_COLORS          = 8
    TELNET_MTTS_MOUSE_TRACKING      = 16
    TELNET_MTTS_OSC_COLOR_PALETTE   = 32
    TELNET_MTTS_SCREEN_READER       = 64
    TELNET_MTTS_PROXY               = 128
    

    # NEW-ENVIRON/ENVIRON codes. 
    TELNET_ENVIRON_IS = 0
    TELNET_ENVIRON_SEND = 1
    TELNET_ENVIRON_INFO = 2
    TELNET_ENVIRON_VAR = 0
    TELNET_ENVIRON_VALUE = 1
    TELNET_ENVIRON_ESC = 2
    TELNET_ENVIRON_USERVAR = 3

    # MSSP codes. 
    TELNET_MSSP_VAR = 1
    TELNET_MSSP_VAL = 2
    
    # MSDP values.
    TELNET_MSDP_VAR         = 1
    TELNET_MSDP_VAL         = 2
    TELNET_MSDP_TABLE_OPEN  = 3
    TELNET_MSDP_TABLE_CLOSE = 4
    TELNET_MSDP_ARRAY_OPEN  = 5
    TELNET_MSDP_ARRAY_CLOSE = 6
    
    # newline, cr and nul
    TELNET_CR = 13
    TELNET_NL = 10
    TELNET_NUL = 0
  end
end
