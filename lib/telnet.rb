
# file::    telnetfilter.rb
# author::  Jon A. Lambert
# version:: 2.8.0
# date::    01/19/2006
#
# This source code copyright (C) 2005, 2006 by Jon A. Lambert
# All rights reserved.
#
# Released under the terms of the TeensyMUD Public License
# See LICENSE file for additional information.
#
$:.unshift "lib" if !$:.include? "lib"
$:.unshift "vendor" if !$:.include? "vendor"

require 'strscan'
require 'ostruct'
require 'network/protocol/filter'
require 'network/protocol/telnetcodes'
require 'network/protocol/asciicodes'

# The TelnetFilter class implements the Telnet protocol.
#
# This implements most of basic Telnet as per RFCs 854/855/1129/1143 and
# options in RFCs 857/858/1073/1091
#
class Telnet
  include Telnet::Codes

  # Initialize state of filter
  def initialize(server)
    @server     = server
    @mode       = :normal #  Parse mode :normal, :cmd, :cr
    @state      = {}
    @sc         = nil
    @sneg_opts  = [ TTYPE, ZMP ]                                
    @ttype      = []
    @init_tries = 0   # Number of tries at negotitating sub options
    @synch      = false
    log.debug "telnet filter initialized - #{@init_tries}"
  end
  
  # Wait for input from the server 
  def wait_for_input
    return Fiber.yield
  end
  
  # Called when client data should be filtered before being passed to the server
  def client_to_server(data)
    
  end
  
  # Called when server data should be filtered before being passed to the client
  def server_to_client(data)
  
  end
  
  
  

  # Negotiate starting wanted options
  #
  # [+args+] Optional initial options
  def init(args)
    if @server.service_type == :client  # let server offer and ask for client
      # several sorts of options here - server offer, ask client or both
      @wopts.each do |key,val|
        case key
        when ECHO, SGA, BINARY, ZMP, EOREC
          ask_him(key,val)
        else
          offer_us(key,val)
        end
      end
    else
      # several sorts of options here - server offer, ask client or both
      @wopts.each do |key,val|
        case key
        when ECHO, SGA, BINARY, ZMP, EOREC
          offer_us(key,val)
        else
          ask_him(key,val)
        end
      end
    end
    true
  end

  # The filter_in method filters input data
  # [+str+]    The string to be processed
  # [+return+] The filtered data
  def filter_in(str)
#    init_subneg
    return "" if str.nil? || str.empty?
    buf = ""

    @sc ? @sc.concat(str) : @sc = StringScanner.new(str)
    while b = @sc.get_byte

      # OOB sync data
      if @pstack.urgent_on || b.getbyte(0) == DM
        log.debug("(#{@pstack.conn.object_id}) Sync mode on")
        @pstack.urgent_on = false
        @synch = true
        break
      end

      case mode?
      when :normal
        case b.getbyte(0)
        when CR
          next if @synch
          set_mode(:cr) if !@pstack.binary_on
        when LF  # LF or LF/CR may be issued by broken mud servers and clients
          next if @synch
          set_mode(:lf) if !@pstack.binary_on
          buf << LF.chr
          echo(CR.chr + LF.chr)
        when IAC
          set_mode(:cmd)
        when NUL  # ignore NULs in stream when in normal mode
          next if @synch
          if @pstack.binary_on
            buf << b
            echo(b)
          else
            log.debug("(#{@pstack.conn.object_id}) unexpected NUL found in stream")
          end
        when BS, DEL
          next if @synch
          # Leaves BS, DEL in input stream for higher filter to deal with.
          buf << b
          echo(BS.chr)
        else
          next if @synch
          ### NOTE - we will allow 8-bit NVT against RFC 1123 recommendation "should not"
          ###
          # Only let 7-bit values through in normal mode
          #if (b[0] & 0x80 == 0) && !@pstack.binary_on
            buf << b
            echo(b)
          #else
          #  log.debug("(#{@pstack.conn.object_id}) unexpected 8-bit byte found in stream '#{b[0]}'")
          #end
        end
      when :cr
        # handle CRLF and CRNUL by insertion of LF into buffer
        case b.getbyte(0)
        when LF
          buf << LF.chr
          echo(CR.chr + LF.chr)
        when NUL
          if @server.service_type == :client  # Don't xlate CRNUL when client
            buf << CR.chr
            echo(CR.chr)
          else
            buf << LF.chr
            echo(CR.chr + LF.chr)
          end
        else # eat lone CR
          buf << b
          echo(b)
        end
        set_mode(:normal)
      when :lf
        # liberally handle LF, LFCR for clients that aren't telnet correct
        case b.getbyte(0)
        when CR # Handle LFCR by swallowing CR
        else  # Handle other stuff that follows - single LF
          buf << b
          echo(b)
        end
        set_mode(:normal)
      when :cmd
        case b.getbyte(0)
        when IAC
          # IAC escapes IAC
          buf << IAC.chr
          set_mode(:normal)
        when AYT
          log.debug("(#{@pstack.conn.object_id}) AYT sent - Msg returned")
          @pstack.conn.sock.send("WOE is here.\n",0)
          set_mode(:normal)
        when AO
          log.debug("(#{@pstack.conn.object_id}) AO sent - Synch returned")
          @pstack.conn.sockio.write_flush
          @pstack.conn.sock.send(IAC.chr + DM.chr, 0)
          @pstack.conn.sockio.write_urgent(DM.chr)
          set_mode(:normal)
        when IP
          @pstack.conn.sockio.read_flush
          @pstack.conn.sockio.write_flush
          log.debug("(#{@pstack.conn.object_id}) IP sent")
          set_mode(:normal)
        when GA, NOP, BRK  # not implemented or ignored
          log.debug("(#{@pstack.conn.object_id}) GA, NOP or BRK sent")
          set_mode(:normal)
        when DM
          log.debug("(#{@pstack.conn.object_id}) Synch mode off")
          @synch = false
          set_mode(:normal)
        when EC
          next if @synch
          log.debug("(#{@pstack.conn.object_id}) EC sent")
          if buf.size > 1
            buf.slice!(-1)
          elsif @pstack.conn.inbuffer.size > 0
            @pstack.conn.inbuffer.slice(-1)
          end
          set_mode(:normal)
        when EL
          next if @synch
          log.debug("(#{@pstack.conn.object_id}) EL sent")
          p = buf.rindex("\n")
          if p
            buf.slice!(p+1..-1)
          else
            buf = ""
            p = @pstack.conn.inbuffer.rindex("\n")
            if p
              @pstack.conn.inbuffer.slice!(p+1..-1)
            end
          end
          set_mode(:normal)
        when DO, DONT, WILL, WONT
          if @sc.eos?
            @sc.unscan
            break
          end
          opt = @sc.get_byte
          case b.getbyte(0)
          when WILL
            replies_him(opt.getbyte(0),true)
          when WONT
            replies_him(opt.getbyte(0),false)
          when DO
            requests_us(opt.getbyte(0),true)
          when DONT
            requests_us(opt.getbyte(0),false)
          end
          # Update interesting things in ProtocolStack after negotiation
          case opt.getbyte(0)
          when ECHO
            @pstack.echo_on = enabled?(ECHO, :us)
          when BINARY
            @pstack.binary_on = enabled?(BINARY, :us)
          when ZMP
            @pstack.zmp_on = enabled?(ZMP, :us)
          end
          set_mode(:normal)
        when SB
          @sc.unscan
          break if @sc.check_until(/#{IAC.chr}#{SE.chr}/).nil?
          @sc.get_byte
          opt = @sc.get_byte
          data = @sc.scan_until(/#{IAC.chr}#{SE.chr}/).chop.chop
          parse_subneg(opt.getbyte(0),data)
          set_mode(:normal)
        else
          log.debug("(#{@pstack.conn.object_id}) Unknown Telnet command - #{b.getbyte(0)}")
          set_mode(:normal)
        end
      end
    end  # while b

    @sc = nil if @sc.eos?
    buf
  end

  # The filter_out method filters output data
  # [+str+]    The string to be processed
  # [+return+] The filtered data
  def filter_out(str)
    return '' if str.nil? || str.empty?
    if !@pstack.binary_on
      str.gsub!(/\n/, "\r\n")
    end
    str
  end

  ###### Custom public methods

  # Test to see if option is enabled
  # [+opt+] The Telnet option code
  # [+who+] The side to check :us or :him
  def enabled?(opt, who)
    option(opt)
    e = @state[opt].send(who)
    e == :yes ? true : false
  end

  # Test to see which state we prefer this option to be in
  # [+opt+] The Telnet option code
  def desired?(opt)
    st = @wopts[opt]
    st = false if st.nil?
    st
  end

  # Handle server-side echo
  # [+ch+] character string to echo
  def echo(ch)
    return if @server.service_type == :client  # Never echo for server when client
                                  # Remove this if it makes sense for peer to peer
    if @pstack.echo_on
      if @pstack.hide_on && ch.getbyte(0) != CR
        @pstack.conn.sock.send('*',0)
      else
        @pstack.conn.sock.send(ch,0)
      end
    end
  end

  # Negotiate starting wanted options that imply subnegotation
  # So far only terminal type
  def init_subneg
    return if @init_tries > 20
    @init_tries += 1
    @wopts.each_key do |opt|
      next if !@sneg_opts.include?(opt)
      log.debug("(#{@pstack.conn.object_id}) Subnegotiation attempt for option #{opt}.")
      case opt
      when TTYPE
        who = :him
      else
        who = :us
      end
      if desired?(opt) == enabled?(opt, who)
        case opt
        when TTYPE
          @pstack.conn.sendmsg(IAC.chr + SB.chr + TTYPE.chr + 1.chr + IAC.chr + SE.chr)
        when ZMP
          log.info("(#{@pstack.conn.object_id}) ZMP successfully negotiated." )
          @pstack.conn.sendmsg("#{IAC.chr}#{SB.chr}#{ZMP.chr}" +
            "zmp.check#{NUL.chr}color.#{NUL.chr}" +
            "#{IAC.chr}#{SE.chr}")
          @pstack.conn.sendmsg("#{IAC.chr}#{SB.chr}#{ZMP.chr}" +
            "zmp.ident#{NUL.chr}WOE#{NUL.chr}#{Version}#{NUL.chr}A mud based on the TeensyMUD server#{NUL.chr}" +
            "#{IAC.chr}#{SE.chr}")
          @pstack.conn.sendmsg("#{IAC.chr}#{SB.chr}#{ZMP.chr}" +
            "zmp.ping#{NUL.chr}" +
            "#{IAC.chr}#{SE.chr}")
          @pstack.conn.sendmsg("#{IAC.chr}#{SB.chr}#{ZMP.chr}" +
            "zmp.input#{NUL.chr}\n     I see you support...\n     ZMP protocol\n#{NUL.chr}" +
            "#{IAC.chr}#{SE.chr}")
        end
        @sneg_opts.delete(opt)
      end
    end

    if @init_tries > 20
      log.debug("(#{@pstack.conn.object_id}) Telnet init_subneg option - Timed out after #{@init_tries} tries.")
      @sneg_opts = []
      @pstack.conn.set_initdone
      if !@pstack.terminal or @pstack.terminal.empty?
        @pstack.terminal = "dumb"
      end
    end
  end

  def send_naws
    return if !enabled?(NAWS, :us)
    ts = @pstack.query(:termsize)
    data = [ts[0]].pack('n') + [ts[1]].pack('n')
    data.gsub!(/#{IAC}/, IAC.chr + IAC.chr) # 255 needs to be doubled
    @pstack.conn.sendmsg(IAC.chr + SB.chr + NAWS.chr + data + IAC.chr + SE.chr)
  end

private
  ###### Private methods

  def getopts(wopts)
    # supported options
    wopts.each do |op|
      case op
      when :ttype
        @wopts[TTYPE] = true
      when :echo
        @wopts[ECHO] = true
      when :sga
        @wopts[SGA] = true
      when :naws
        @wopts[NAWS] = true
      when :eorec
        @wopts[EOREC] = true
      when :binary
        @wopts[BINARY] = true
      when :zmp
        @wopts[ZMP] = true
      end
    end
  end

  # parse the subnegotiation data and save it
  # [+opt+] The Telnet option found
  # [+data+] The data found between SB OPTION and IAC SE
  def parse_subneg(opt,data)
    data.gsub!(/#{IAC}#{IAC}/, IAC.chr) # 255 needs to be undoubled from all data
    case opt
    when NAWS
      @pstack.twidth = data[0..1].unpack('n')[0]
      @pstack.theight = data[2..3].unpack('n')[0]
      @pstack.conn.publish(:termsize)
      log.debug("(#{@pstack.conn.object_id}) Terminal width #{@pstack.twidth} / height #{@pstack.theight}")
    when TTYPE
      if data.getbyte(0) == 0
        log.debug("(#{@pstack.conn.object_id}) Terminal type - #{data[1..-1]}")
        if !@ttype.include?(data[1..-1])
          # short-circuit choice because of Zmud
          if data[1..-1].downcase == 'zmud'
            @ttype << data[1..-1]
            @pstack.terminal = 'zmud'
            log.debug("(#{@pstack.conn.object_id}) Terminal choice - #{@pstack.terminal} in list #{@ttype.inspect}")
          end
          # short-circuit choice because of Windows telnet client
          if data[1..-1].downcase == 'vt100'
            @ttype << data[1..-1]
            @pstack.terminal = 'vt100'
            log.debug("(#{@pstack.conn.object_id}) Terminal choice - #{@pstack.terminal} in list #{@ttype.inspect}")
          end
          return if @pstack.terminal
          @ttype << data[1..-1]
          @pstack.conn.sendmsg(IAC.chr + SB.chr + TTYPE.chr + 1.chr + IAC.chr + SE.chr)
        else
          return if @pstack.terminal
          choose_terminal
        end
      elsif data.getbyte(0) == 1  # send - should only be called by :client
        return if !@pstack.terminal
        @pstack.conn.sendmsg(IAC.chr + SB.chr + TTYPE.chr + 0.chr + @pstack.terminal + IAC.chr + SE.chr)
      end
    when ZMP
      args = data.split("\0")
      cmd = args.shift
      handle_zmp(cmd,args)
    end
  end

  # Pick a preferred terminal
  # Order is vt100, vt999, ansi, xterm, or a recognized custom client
  # Should not pick vtnt as we dont handle it
  def choose_terminal
    if @ttype.empty?
      @pstack.terminal = "dumb"
    end

    # Pick most capable from list of terminals
    @pstack.terminal = @ttype.find {|t| t =~ /mushclient/i } if !@pstack.terminal
    @pstack.terminal = @ttype.find {|t| t =~ /simplemu/i } if !@pstack.terminal
    @pstack.terminal = @ttype.find {|t| t =~ /(zmud).*/i } if !@pstack.terminal
    @pstack.terminal = @ttype.find {|t| t =~ /linux/i } if !@pstack.terminal
    @pstack.terminal = @ttype.find {|t| t =~ /cygwin/i } if !@pstack.terminal
    @pstack.terminal = @ttype.find {|t| t =~ /(cons25).*/i } if !@pstack.terminal
    @pstack.terminal = @ttype.find {|t| t =~ /(xterm).*/i } if !@pstack.terminal
    @pstack.terminal = @ttype.find {|t| t =~  /(vt)[-]?100/i } if !@pstack.terminal
    @pstack.terminal = @ttype.find {|t| t =~ /(vt)[-]?\d+/i } if !@pstack.terminal
    @pstack.terminal = @ttype.find {|t| t =~ /(ansi).*/i } if !@pstack.terminal

    if @pstack.terminal && @ttype.last != @pstack.terminal # short circuit retraversal of options
      @ttype.each do |t|
        @pstack.conn.sendmsg(IAC.chr + SB.chr + TTYPE.chr + 1.chr + IAC.chr + SE.chr)
        break if t == @pstack.terminal
      end
    elsif @ttype.last != @pstack.terminal
      @pstack.terminal = 'dumb'
    end

    @pstack.terminal.downcase!

    # translate certain terminals to something meaningful
    case @pstack.terminal
    when /cygwin/i, /cons25/i, /linux/i, /dec-vt/i
      @pstack.terminal = 'vt100'
    when /ansis/i then
      @pstack.terminal = 'ansi'
    end
    log.debug("(#{@pstack.conn.object_id}) Terminal set to - #{@pstack.terminal} from list #{@ttype.inspect}")
  end

  # Get current parse mode
  # [+return+] The current parse mode
  def mode?
    return @mode
  end

  # set current parse mode
  # [+m+] Mode to set it to
  def set_mode(m)
    @mode = m
  end

  # Creates an option entry in our state table and sets its initial state
  def option(opt)
    return if @state.key?(opt)
    o = OpenStruct.new
    o.us = :no
    o.him = :no
    o.usq = :empty
    o.himq = :empty
    @state[opt] = o
  end

  # Ask the client to enable or disable an option.
  #
  # [+opt+]   The option code
  # [+enable+] true for enable, false for disable
  def ask_him(opt, enable)
    log.debug("(#{@pstack.conn.object_id}) Requested Telnet option #{opt.to_s} set to #{enable.to_s}")
    initiate(opt, enable, :him)
  end

  # Offer the server to enable or disable an option
  #
  # [+opt+]   The option code
  # [+enable+] true for enable, false for disable
  def offer_us(opt, enable)
    log.debug("(#{@pstack.conn.object_id}) Offered Telnet option #{opt.to_s} set to #{enable.to_s}")
    initiate(opt, enable, :us)
  end

  # Initiate a request to client.  Called by ask_him or offer_us.
  #
  # [+opt+]   The option code
  # [+enable+] true for enable, false for disable
  # [+who+] :him if asking client, :us if server offering
  def initiate(opt, enable, who)
    option(opt)

    case who
    when :him
      willdo = DO.chr
      wontdont = DONT.chr
      whoq = :himq
    when :us
      willdo = WILL.chr
      wontdont = WONT.chr
      whoq = :usq
    else
      # Error
    end

    case @state[opt].send(who)
    when :no
      if enable
        @state[opt].send("#{who}=", :wantyes)
        @pstack.conn.sendmsg(IAC.chr + willdo + opt.chr)
      else
        # Error already disabled
        log.error("(#{@pstack.conn.object_id}) Telnet negotiation: option #{opt.to_s} already disabled")
      end
    when :yes
      if enable
        # Error already enabled
        log.error("(#{@pstack.conn.object_id}) Telnet negotiation: option #{opt.to_s} already enabled")
      else
        @state[opt].send("#{who}=", :wantno)
        @pstack.conn.sendmsg(IAC.chr + wontdont + opt.chr)
      end
    when :wantno
      if enable
        case @state[opt].send(whoq)
        when :empty
          @state[opt].send("#{whoq}=", :opposite)
        when :opposite
          # Error already queued enable request
          log.error("(#{@pstack.conn.object_id}) Telnet negotiation: option #{opt.to_s} already queued enable request")
        end
      else
        case @state[opt].send(whoq)
        when :empty
          # Error already negotiating for disable
          log.error("(#{@pstack.conn.object_id}) Telnet negotiation: option #{opt.to_s} already negotiating for disable")
        when :opposite
          @state[opt].send("#{whoq}=", :empty)
        end
      end
    when :wantyes
      if enable
        case @state[opt].send(whoq)
        when :empty
          #Error already negotiating for enable
          log.error("(#{@pstack.conn.object_id}) Telnet negotiation: option #{opt.to_s} already negotiating for enable")
        when :opposite
          @state[opt].send("#{whoq}=", :empty)
        end
      else
        case @state[opt].send(whoq)
        when :empty
          @state[opt].send("#{whoq}=", :opposite)
        when :opposite
          #Error already queued for disable request
          log.error("(#{@pstack.conn.object_id}) Telnet negotiation: option #{opt.to_s} already queued for disable request")
        end
      end
    end
  end

  # Client replies WILL or WONT
  #
  # [+opt+]   The option code
  # [+enable+] true for WILL answer, false for WONT answer
  def replies_him(opt, enable)
    log.debug("(#{@pstack.conn.object_id}) Client replies to Telnet option #{opt.to_s} set to #{enable.to_s}")
    response(opt, enable, :him)
  end

  # Client requests DO or DONT
  #
  # [+opt+]   The option code
  # [+enable+] true for DO request, false for DONT request
  def requests_us(opt, enable)
    log.debug("(#{@pstack.conn.object_id}) Client requests Telnet option #{opt.to_s} set to #{enable.to_s}")
    response(opt, enable, :us)
  end

  # Handle client response.  Called by requests_us or replies_him
  #
  # [+opt+]   The option code
  # [+enable+] true for WILL answer, false for WONT answer
  # [+who+] :him if client replies, :us if client requests
  def response(opt, enable, who)
    option(opt)

    case who
    when :him
      willdo = DO.chr
      wontdont = DONT.chr
      whoq = :himq
    when :us
      willdo = WILL.chr
      wontdont = WONT.chr
      whoq = :usq
    else
      # Error
    end

    case @state[opt].send(who)
    when :no
      if enable
        if desired?(opt)
        # If we agree
          @state[opt].send("#{who}=", :yes)
          @pstack.conn.sendmsg(IAC.chr + willdo + opt.chr)
          log.debug("(#{@pstack.conn.object_id}) Telnet negotiation: agreed to enable option #{opt}")
        else
        # If we disagree
          @pstack.conn.sendmsg(IAC.chr + wontdont + opt.chr)
          log.debug("(#{@pstack.conn.object_id}) Telnet negotiation: disagreed to enable option #{opt}")
        end
      else
        # Ignore
      end
    when :yes
      if enable
        # Ignore
      else
        @state[opt].send("#{who}=", :no)
        @pstack.conn.sendmsg(IAC.chr + wontdont + opt.chr)
      end
    when :wantno
      if enable
        case @state[opt].send(whoq)
        when :empty
          #Error DONT/WONT answered by WILL/DO
          @state[opt].send("#{who}=", :no)
        when :opposite
          #Error DONT/WONT answered by WILL/DO
          @state[opt].send("#{who}=", :yes)
          @state[opt].send("#{whoq}=", :empty)
        end
        log.error("(#{@pstack.conn.object_id}) Telnet negotiation: option #{opt.to_s} DONT/WONT answered by WILL/DO")
      else
        case @state[opt].send(whoq)
        when :empty
          @state[opt].send("#{who}=", :no)
          log.debug("(#{@pstack.conn.object_id}) Telnet negotiation: agreed to disable option #{opt}")
        when :opposite
          @state[opt].send("#{who}=", :wantyes)
          @state[opt].send("#{whoq}=", :empty)
          @pstack.conn.sendmsg(IAC.chr + willdo + opt.chr)
        end
      end
    when :wantyes
      if enable
        case @state[opt].send(whoq)
        when :empty
          @state[opt].send("#{who}=", :yes)
          log.debug("(#{@pstack.conn.object_id}) Telnet negotiation: agreed to enable option #{opt}")
        when :opposite
          @state[opt].send("#{who}=", :wantno)
          @state[opt].send("#{whoq}=", :empty)
          @pstack.conn.sendmsg(IAC.chr + wontdont + opt.chr)
        end
      else
        case @state[opt].send(whoq)
        when :empty
          @state[opt].send("#{who}=", :no)
          log.debug("(#{@pstack.conn.object_id}) Telnet negotiation: agreed to disable option #{opt}")
        when :opposite
          @state[opt].send("#{who}=", :no)
          @state[opt].send("#{whoq}=", :empty)
        end
      end
    end
  end

  def handle_zmp(cmd,args)
    log.debug("(#{@pstack.conn.object_id}) ZMP command recieved - '#{cmd}' args: #{args.inspect}" )
    case cmd
    when "zmp.ping"
      @pstack.conn.sendmsg("#{IAC.chr}#{SB.chr}#{ZMP.chr}" +
        "zmp.time#{NUL.chr}#{Time.now.utc.strftime("%Y-%m-%d %H:%M:%S")}#{NUL.chr}" +
        "#{IAC.chr}#{SE.chr}")
    when "zmp.time"
    when "zmp.ident"
      # That's nice
    when "zmp.check"
      case args[0]
      when /zmp.*/
      # We support all 'zmp.' package and commands so..
        @pstack.conn.sendmsg("#{IAC.chr}#{SB.chr}#{ZMP.chr}" +
          "zmp.support#{NUL.chr}#{args[0]}{NUL.chr}" +
          "#{IAC.chr}#{SE.chr}")
      else
        @pstack.conn.sendmsg("#{IAC.chr}#{SB.chr}#{ZMP.chr}" +
          "zmp.no-support#{NUL.chr}#{args[0]}#{NUL.chr}" +
          "#{IAC.chr}#{SE.chr}")
      end
    when "zmp.support"
    when "zmp.no-support"
    when "zmp.input"
      # Now we just simply pass this whole load to the Character.parse
      # WARN: This means there is a possibility of out-of-order processing
      #       of @inbuffer, though extremely unlikely.
      @pstack.conn.publish(args[0])
    end
  end

end
