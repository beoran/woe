
class Dir
  def self.mkdir_p(name)
    sub   = ""
    parts = name.split('/').reject { |e| e.empty? }
    parts.each do | part |
      sub <<  "/#{part}"
      mkdir sub
    end
  end
end



# Module to help with serialization and deserialization of any type of data
module Serdes
  
  module ClassMethods
    def serdes_add_to_register(name)
      @serder_register ||= []
      @serder_register << name
    end
    
    def serdes_reader(name)
      serdes_add_to_register(name)
      attr_reader(name)
    end
    
    def serdes_writer(name)
      serdes_add_to_register(name)
      attr_writer(name)
    end
    
    def serdes_accessor(name)
      serdes_add_to_register(name)
      attr_accessor(name)
    end
    
    def serdes_register()
      @serder_register ||= []
      return @serder_register
    end  
  end
  
  # include callback, be sure to extend the class with the ClassMethods
  def self.included(klass)
    klass.extend(ClassMethods)
  end
  
  def serdes_dir_name
    top = self.class.to_s.gsub('::', '/').downcase
    top << '/' 
    top << self.id.to_s[0]
    top << '/' 
    top << self.id.to_s    
    return top
  end
  
  def serdes_file_name
    top = self.id.to_s.dup    
    top << '.'
    top << self.class.to_s.gsub('::', '.').downcase
    return top 
  end

  
  def save
    Dir.mkdir_p serdes_dir_name
    data = {}
    self.class.serdes_register.each do |name|
      key = "#{name}" 
      value = "#{self.send(name.to_sym)}"
      data[key] = value
    end
    full_name = serdes_dir_name + '/' + serdes_file_name
    Sitef.save_filename(full_name, [ data ] )
  end
  


end
