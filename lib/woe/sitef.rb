
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
  # Keys may not be nested, however, une could use spaces or dots, 
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
      # XXX does eof? even work???
      break if line.nil?
      next if line.empty? 
      # new record
      if line[0,2] == '--' 
        # Store last key used if any.
        if key
          record[key.downcase] = value.chomp
          key = nil
        end  
        results << record
        record = {}
      elsif line[0] == '#'
      # Comments start with #
      elsif line[0] == ':'
      # a key/value pair
      key, value = line[1,line.size].split(':', 2)
      record[key.downcase] = value.chomp
      key = value = nil
      elsif line[0, 2] == '..'
      # end of multiline value 
      record[key.downcase] = value.chomp
      key = value = nil
      elsif (line[0] == '.') && key.nil?
      # Multiline key/value starts here (but is ignored 
      # until .. is encountered)
      key   = line[1, line.size]
      value = ""
      elsif key
          # continue the value
          value << line
      else
          # Not in a key, sntax error.
          errors << "#{lineno}: Don't know how to process line"
      end      
    end
    # Store last key used if any.
    if key
      record[key.downcase] = value.chomp
    end  
    # store last record 
    results << record unless record.empty?
    return results, errors
  end  
  
  def self.load_filename(filename)
    results , errors, warnings = nil, nil, nil;
    file = File.open(filename, 'rt')
    return nil, ["Could not open #{filename}"] unless file
    begin 
      results, errors = parse_file(file)
    ensure
      file.close
    end
    return results, errors
  end
  
  def self.save_field(file, key, value)
    sval = value.to_s
    if sval["\n"]
      file.puts(".#{key}\n")
      file.puts(sval)
      file.puts("\n..\n")
    else
      file.puts(":#{key}:#{sval}\n")
    end
  end
  
  def self.save_record(file, record)
    record.each do | key, value |
      save_field(file, key, value)
    end
  end

  def self.save_file(file, records)
    records.each do | record |
      save_record(file, record)
      file.puts("--\n")
    end
  end
  
  def self.save_filename(filename, records)
    results , errors = nil, nil
    file = File.open(filename, 'wt')
    return false, ["Could not open #{filename}"] unless file
    begin 
      save_file(file, records)
    ensure
      file.close
    end
    return true, []
  end
  
end

