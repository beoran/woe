
# Sitef is a simple text format for serializing data to
# It's intent is to be human readable and easy to 
# use for multi line text.

module Sitef
  # All Sitef data is stored in files with one or more records.
  # Records are separated by separated by at least 2 dashes on a line.
  # Records contain key/value fields. The key starts in the first column
  # with a : and is followed by a : the value starts after the second :
  # A multiline key / value needs a key that starts with . and ends with .
  # the end of the value is a  pair of dots .. by itself 
  # Keys may not be nested, however, you could use spaces or dots, 
  # or array indexes to emulate nexted keys. 
  # A # at the start optionally after whitespace is a comment
  # 
  def self.parse_file(file)
    lineno   = 0
    results  = []
    errors   = []
    
    record   = {}
    key      = nil
    value    = nil
    until file.eof?
      lineno     += 1 
      line        = file.gets(256)
      break if line.nil?
      next if line.empty? 
      # new record
      if line[0,2] == '--' 
        # Store last key used if any.
        if key          
          record[key] = value.chomp
          key = nil
        end  
        results << record
        record = {}
      elsif line[0] == '#'
      # Comments start with #
      elsif line[0] == ':'
      # a key/value pair
      key, value  = line[1,line.size].split(':', 2)
      record[key] = value.chomp
      key = value = nil
      elsif line[0, 2] == '..'
      # end of multiline value 
      record[key] = value.chomp
      key = value = nil
      elsif (line[0] == '.') && key.nil?
      # Multiline key/value starts here (but is ignored 
      # until .. is encountered)
      key   = line[1, line.size]
      key.chomp!
      value = ""
      # multiline value
      elsif key
          if line[0] == '\\'
            # remove any escapes
            line.slice!(0)
          end
          # continue the value
          value << line
      else
          # Not in a key, sntax error.
          errors << "#{lineno}: Don't know how to process line"
      end      
    end
    # Store last key used if any.
    if key      
      record[key] = value.chomp
    end  
    # store last record 
    results << record unless record.empty?
    return results, errors
  end  
  
  def self.load_filename(filename)
    results, errors = nil, nil, nil;
    file = File.open(filename, 'rt')
    return nil, ["Could not open #{filename}"] unless file
    begin 
      results, errors = parse_file(file)
    ensure
      file.close
    end
    return results, errors
  end
  
  # Loads a Sitef fileas obejcts. Uses the ruby_klass atribute to load the object
  # If that is missing, uses defklass
  def self.load_objects(filename, defklass=nil)
    results, errors = load_filename(filename)
    p filename, results, errors
    unless errors.nil? || errors.empty?
      return nil, errors 
    end
    
    objres = [] 
    results.each do | result |
      klassname = result['ruby_class'] || defklass
      return nil unless klassname
      klass = klassname.split('::').inject(Kernel) { |klass, name| klass.const_get(name) rescue nil } 
      return nil unless klass
      if klass.respond_to? :from_sitef
        objres << klass.from_sitef(result)
      else
        objres << klass.new(result)
      end      
    end
    return objres, errors    
  end
  
  
  # Saves a single field to a file in Sitef format.
  def self.save_field(file, key, value)
    if value.is_a? String
      sval = value.dup
    else
      sval = value.to_s
    end
    if sval["\n"]
      file.puts(".#{key}\n")
      # Escape everything that could be misinterpreted with a \\
      sval.gsub!(/\n([\.\-\:\#\\]+)/, "\n\\\\\\1")
      sval.gsub!(/\A([\.\-\:\#\\]+)/, "\\\\\\1")
      file.printf("%s", sval)
      file.printf("\n..\n")
    else
      file.printf(":#{key}:#{sval}\n")
    end
  end
  
  def self.save_object(file, object, *fields)
    save_field(file, :ruby_class, object.class.to_s)
    fields.each do | field |
      value = object.send(field.to_sym)
      save_field(file, field, value)
    end
  end
  
  def self.save_record(file, record, *fields)
    record.each do | key, value |
      next if fields && !fields.empty? && !fields.member?(key)
      save_field(file, key, value)
    end
  end

  def self.save_file(file, records, *fields)
    records.each do | record |
      if record.is_a? Hash
        save_record(file, record, *fields)
      else 
        save_object(file, record, *fields)
      end
      file.puts("--\n")
    end
  end
  
  def self.save_filename(filename, records, *fields)
    results , errors = nil, nil
    file = File.open(filename, 'wt')
    return false, ["Could not open #{filename}"] unless file
    begin 
      save_file(file, records, *fields)
    ensure
      file.close
    end
    return true, []
  end
  
end

