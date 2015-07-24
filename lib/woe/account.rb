require_relative "../sitef.rb"
require_relative "../serdes.rb"
require_relative "../monolog.rb"
require_relative "../security.rb"



module Woe
class Account
  include Serdes
  include Monolog

  serdes_reader   :id
  serdes_reader   :pass
  serdes_reader   :algo
  serdes_reader   :email
  serdes_accessor :woe_points
  
  
  def inspect
    "Account #{@id} #{@pass} #{algo}"
  end
  
  def password=(pass)
    @algo = "crypt"
    @pass = crypt(pass)
  end
  
  # Returns true if the password matches that of this account or false if not.
  def challenge?(trypass) 
    if algo == "plain"
      return @pass == trypass
    elsif algo == "crypt"
      return crypt_challenge?(trypass, @pass)
    else
      return false
    end
  end
end # class Account
end # module Woe






