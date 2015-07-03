script "serdes.rb"


class Account
  include Serdes

  serdes_reader :id
  serdes_reader :pass
  serdes_reader :algo
  
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
      return crypt_challenge(trypass, @pass)
    else
      return false
    end
  end
  
end






