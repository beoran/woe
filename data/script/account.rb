script "serdes.rb"


class Account
  include Serdes
  
  serdes_reader :id
  serdes_reader :pass
  serdes_reader :algo
  
  def initialize(nam, pas, alg = 'plain') 
    @id   = nam
    @algo = alg
    @pass = pas
  end

end






