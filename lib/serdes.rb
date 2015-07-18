


class Dir
  def self.mkdir_p(name)
    sub   = ""
    parts = name.split('/').reject { |e| e.empty? }
    parts.each do | part |
      sub <<  "/#{part}"
      mkdir sub rescue nil
    end
  end
end



# Module to help with serialization and deserialization of any type of data
module Serdes
  
  module ClassMethods
    def serdes_add_to_fields(name, type = nil)
      @serdes_fields ||= []
      info = { :name => name, :type => type }
      @serdes_fields << info
    end
    
    def serdes_reader(name, type = nil)
      serdes_add_to_fields(name, type)
      attr_reader(name)
    end
    
    def serdes_writer(name)
      serdes_add_to_fields(name, type = nil)
      attr_writer(name)
    end
    
    def serdes_accessor(name, type = nil)
      serdes_add_to_fields(name, type)
      attr_accessor(name)
    end
    
    def serdes_fields()
      @serdes_fields ||= []
      return @serdes_fields
    end
    
    
    def serdes_register(obj)
      @serdes_loaded ||= {}
      @serdes_loaded[obj.id] = obj
    end
    
    def serdes_forget(id)
      @serdes_loaded ||= {}
      @serdes_loaded.delete(id)
    end
    
    def serdes_loaded()
      @serdes_loaded ||= {}
      return @serdes_loaded
    end
    
    def serdes_get(id)
      @serdes_loaded ||= {}
      return @serdes_loaded[id.to_sym]
    end

    def serdes_load_one(id)
      return nil unless id && !id.empty?
      
      full_name = Serdes.serdes_full_for(self, id)
      data, errors  = Sitef.load_filename(full_name)
      unless data
        # log_error(errors.join(", "))
        return nil
      end
      
      eldat = data.select do |el|
        el['id'].to_s == id.to_s
      end
      return nil unless eldat

      eldat = eldat.first
      return nil unless eldat
      
      typedat = {}
      self.serdes_fields.each do |info|
        name  = info[:name]
        type  = info[:type]
        value = eldat[name.to_s]
        
        typevalue = nil
        
        if type
          if type.respond_to?(:serdes_load)
            typevalue = type.serdes_load(value)
          elsif Kernel.respond_to?(type.to_sym)
            typevalue = Kernel.send(type.to_sym, value) rescue nil 
          elsif type.respond_to(:new)
            typevalue = type.new(value)
          else 
            typevalue = value
          end
        else
          typevalue = value
        end
      
        typedat[name] = typevalue
      end
      
      obj = self.new(typedat)
      return obj
    end
    
    def serdes_fetch(id)
      res = serdes_get(id)
      return res if res
      return serdes_load_one(id)
    end
    
    alias :fetch    :serdes_fetch
    alias :load_one :serdes_load_one
    alias :get      :serdes_get
    
    def from_serdes(id)
      return serdes_fetch(id)
    end
    
    def to_serdes(value)
      return value.id.to_s
    end
    
    
  end

  # include callback, be sure to extend the class with the ClassMethods
  def self.included(klass)
    klass.extend(ClassMethods)
  end
  
  def self.serdes_dir=(dir)
    @serdes_dir = dir
  end
  
  def self.serdes_dir
    @serdes_dir ||= File.join(Dir.pwd, 'data', 'var')
    @serdes_dir
  end

  
  def self.serdes_dir_for(klass)
    top = File.join(Serdes.serdes_dir,
      klass.to_s.gsub('::', '/').downcase)
    return top
  end
  
  def self.serdes_file_for(id)
    top = id.to_s.dup    
    top << '.sitef'
    return top 
  end
  
  def self.serdes_full_for(klass, id)
    tdir = serdes_dir_for(klass)
    tfil = serdes_file_for(id)
    return File.join(tdir, tfil)
  end
  

  def serdes_data
    data = {}
    self.class.serdes_fields.each do |info|
      name  = info[:name]
      type  = info[:type]
      type||= String
      key   = "#{name}" 
      value = "#{self.send(name.to_sym)}"
      if type.respond_to?(:to_serdes)
         wrapvalue = type.to_serdes(value)
      else 
         wrapvalue = value.to_s
      end
      data[key]    = wrapvalue
    end
    return data
  end
  
  def save_one
    Dir.mkdir_p Serdes.serdes_dir_for(self.class)
    data = serdes_data
    full_name = Serdes.serdes_full_for(self.class, self.id)
    Sitef.save_filename(full_name, [ data ] )
  end
  
  def initialize(fields = {}) 
    fields.each  do |key, value|
      p "Setting #{key} #{value}"
      instance_variable_set("@#{key}", value)
    end
    self.class.serdes_register(self)
  end

end

