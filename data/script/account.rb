script "serdes.rb"


class Account
  include Serdes
  
  serdes_reader :id
  serdes_reader :pass
  serdes_reader :algo
  
end






