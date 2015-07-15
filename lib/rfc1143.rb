
require_relative 'telnet/codes'


# rfc1143 state machine for telnet protocol
# Thehandle_ functions handle input,
# the send_ functions are for sending negotiations
class RFC1143
  include Telnet::Codes

  attr_reader :telopt
  attr_reader :us
  attr_reader :him
  attr_reader :agree
  
  def initialize(to, u, h, a)
    @telopt = to
    @us     = u
    @him    = h
    @agree  = a
  end

=begin
 EXAMPLE STATE MACHINE
 FOR THE Q METHOD OF IMPLEMENTING TELNET OPTION NEGOTIATION

    There are two sides, we (us) and he (him).  We keep four
    variables:

       us: state of option on our side (NO/WANTNO/WANTYES/YES)
       usq: a queue bit (EMPTY/OPPOSITE) if us is WANTNO or WANTYES
       him: state of option on his side
       himq: a queue bit if him is WANTNO or WANTYES

    An option is enabled if and only if its state is YES.  Note that
    us/usq and him/himq could be combined into two six-choice states.

    "Error" below means that producing diagnostic information may be a
    good idea, though it isn't required.

    Upon receipt of WILL, we choose based upon him and himq:
       NO            If we agree that he should enable, him=YES, send
                     DO; otherwise, send DONT.
       YES           Ignore.
       WANTNO  EMPTY Error: DONT answered by WILL. him=NO.
            OPPOSITE Error: DONT answered by WILL. him=YES*,
                     himq=EMPTY.
       WANTYES EMPTY him=YES.
            OPPOSITE him=WANTNO, himq=EMPTY, send DONT.

    * This behavior is debatable; DONT will never be answered by WILL
      over a reliable connection between TELNETs compliant with this
      RFC, so this was chosen (1) not to generate further messages,
      because if we know we're dealing with a noncompliant TELNET we
      shouldn't trust it to be sensible; (2) to empty the queue
      sensibly.

=end
  def handle_will
    case @us
    when :no
      if @agree
        return TELNET_DO, @telopt
      else
        return TELNET_DONT, @telopt
      end
    when :yes
      # ignore
      return nil, nil
    when :wantno
      @him = :no
      return :error, "DONT answered by WILL"
    when :wantno_opposite
      @him = :yes
      return :error, "DONT answered by WILL"
    when :wantyes
      @him = :yes
      return nil, nil
    when :wantyes_opposite
      @him = :wantno
      return TELNET_DONT, @telopt
    end
  end
  
  
=begin
Upon receipt of WONT, we choose based upon him and himq:
   NO            Ignore.
   YES           him=NO, send DONT.
   WANTNO  EMPTY him=NO.
        OPPOSITE him=WANTYES, himq=NONE, send DO.
   WANTYES EMPTY him=NO.*
        OPPOSITE him=NO, himq=NONE.**

* Here is the only spot a length-two queue could be useful; after
  a WILL negotiation was refused, a queue of WONT WILL would mean
  to request the option again. This seems of too little utility
  and too much potential waste; there is little chance that the
  other side will change its mind immediately.

** Here we don't have to generate another request because we've
   been "refused into" the correct state anyway.
=end 
  def handle_wont
    case @us
    when :no
      return nil, nil
    when :yes
      @him = :no
      return TELNET_DONT, @telopt
    when :wantno
      @him = :no
      return nil, nil
    when :wantno_opposite
      @him = :wantyes
      return TELNET_DO, @telopt
    when :wantyes
      @him = :no
      return nil, nil
    when :wantyes_opposite
      @him = :no
      return nil, nil
    end
  end   
  
=begin

    If we decide to ask him to enable:
       NO            him=WANTYES, send DO.
       YES           Error: Already enabled.
       WANTNO  EMPTY If we are queueing requests, himq=OPPOSITE;
                     otherwise, Error: Cannot initiate new request
                     in the middle of negotiation.
            OPPOSITE Error: Already queued an enable request.
       WANTYES EMPTY Error: Already negotiating for enable.
            OPPOSITE himq=EMPTY.
            
    We handle the option on our side by the same procedures, with DO-
    WILL, DONT-WONT, him-us, himq-usq swapped.         
=end
  def handle_do
    case @him
    when :no
      @us = :wantyes
      return TELNET_WILL, @telopt
    when :yes
      return :error, 'Already enabled'
    when :wantno
      # us = :wantno_opposite # only if "buffering", whatever that means.
      return :error, 'Request in the middle of negotiation'
    when :wantno_opposite
      return :error, 'Already queued request'
    when :wantyes
      return :error, 'Already negotiating for enable'
    when :wantyes_opposite
      @us = :wantyes
      return nil, nil
    end
  end   
  
=begin
    If we decide to ask him to disable:
       NO            Error: Already disabled.
       YES           him=WANTNO, send DONT.
       WANTNO  EMPTY Error: Already negotiating for disable.
            OPPOSITE himq=EMPTY.
       WANTYES EMPTY If we are queueing requests, himq=OPPOSITE;
                     otherwise, Error: Cannot initiate new request
                     in the middle of negotiation.
            OPPOSITE Error: Already queued a disable request.

    We handle the option on our side by the same procedures, with DO-
    WILL, DONT-WONT, him-us, himq-usq swapped.
=end
  def handle_dont
    case @him
    when :no
      return :error, 'Already disabled'
    when :yes
      @us = :wantno
      return TELNET_WONT, @telopt
    when :wantno
      return :error, 'Already negotiating for disable'
    when :wantno_opposite
      @us = :wantno
      return nil, nil
    when :wantyes
      # us = :wantno_opposite # only if "buffering", whatever that means.
      return :error, 'Request in the middle of negotiation'
    when :wantyes_opposite
      return :error, 'Already queued disable request'
    end
  end
  
  # advertise willingess to support an option 
  def send_will
    case @us
    when :no
      @us = :wantyes
      return TELNET_WILL, @telopt
    when :wantno
      @us = :wantno_opposite
    when :wantyes_opposite
      @us = :wantyes
    else
      return nil, nil
    end
  end
  
  # force turn-off of locally enabled option
  def send_wont
    case @us
    when :yes
      @us = :wantno
      return TELNET_WONT, @telopt
    when :wantno_opposite
      @us = :wantno
      return nil, nil
    when :wantyes
      @us = :wantyes_opposite
      return nil, nil
    else
      return nil, nil
    end
  end   

  # ask remote end to enable an option
  def send_do
    case @him
    when :no
      @him = :wantyes
      return TELNET_DO, @telopt
    when :wantno
      @him = :wantno_opposite
      return nil, nil
    when :wantyes_opposite
      @us = :wantyes
      return nil, nil
    else
      return nil, nil
    end
  end


  # demand remote end disable an option
  def send_dont
    case @him
    when :yes
      @him = :wantno
      return TELNET_DONT, @telopt
    when :wantno_opposite
      @him = :wantno
      return nil, nil
    when :wantyes
      @him = :wantyes_opposite
    else
      return nil, nil
    end
  end
end
  



