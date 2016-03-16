package sitef

import "os"
import "io"
import "strings"
import "fmt"
import "bytes"
import "bufio"
import "strconv"
import "reflect"
import "errors"
import "github.com/beoran/woe/monolog"


// Sitef format for serialization
// Sitef is a simple text format for serializing data to
// It's intent is to be human readable and easy to 
// use for multi line text.
// It is quite similar to recfiles, though not compatible because 
// in sitef files the first character on the line determines the meaning, 
// and there is no backslash escaping.
//
// Sitef is a line based syntax where the first character on the line
// determines the meaning of the line.
// Several lines together form a record.
// A line that starts with # is a comment. There may be no whitespace
// in front of the comment.
// A newline character by itself (that is, an empty line), 
// or a - ends a record.
// A plus character, an escape on the previous line or a tab or 
// a space continues a value.
// A Continues value gets a newline inserted only when a space or tab was used.
// + supresses the newline.
// Anything else signifies the beginning of the next key.
// % is allowed for special keys for recfile compatibility.  
// However % directives are not implemented.
// Keys may not be nested, however, you could use spaces or dots, 
// or array indexes to emulate nexted keys. 
// A # at the start optionally after whitespace is a comment
// 

type Record map[string]string

func (me * Record) Put(key string, val string) {
    (*me)[key] = val
}

func (me * Record) Putf(key string, format string, values ...interface{}) {
    me.Put(key, fmt.Sprintf(format, values...)) 
}

func (me * Record) PutArray(key string, values []string) {
    for i, value := range values {
        realkey := fmt.Sprintf("%s[%d]", key, i)
        me.Put(realkey, value)
    }
} 

func (me * Record) PutInt(key string, val int) {
    me.Putf(key, "%d", val)
}


func (me * Record) PutInt64(key string, val int64) {
    me.Putf(key, "%d", val)
}

func (me * Record) PutFloat64(key string, val float64) {
    me.Putf(key, "%lf", val)
}

func (me Record) MayGet(key string) (result string, ok bool) {
    result, ok = me[key]
    return result, ok
}


func (me Record) Get(key string) (result string) {
    result= me[key]
    return result
}

func (me Record) Getf(key string, format string, 
    values ...interface{}) (amount int, ok bool) {
    val := me.Get(key)
    count, err := fmt.Sscanf(val, format, values...)
    if err != nil {
        return 0, false
    }
    return count, true
}

func (me Record) GetInt(key string) (val int, err error) {
    i, err := strconv.ParseInt(me.Get(key), 0, 0)
    return int(i), err
}

func (me Record) GetIntDefault(key string, def int) (val int) {
    i, err := strconv.ParseInt(me.Get(key), 0, 0)
    if err != nil {
        return def;
    }
    return int(i);
}


func (me Record) GetFloat(key string) (val float64, error error) {
    return strconv.ParseFloat(me.Get(key), 64)
}

func (me * Record) convSimple(typ reflect.Type, val reflect.Value) (res string, err error) {
    switch val.Kind() {
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
    return strconv.FormatInt(val.Int(), 10), nil
        
    case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
    return strconv.FormatUint(val.Uint(), 10), nil
    
    case reflect.Float32, reflect.Float64:
    return strconv.FormatFloat(val.Float(), 'g', -1, val.Type().Bits()), nil
    case reflect.String:
    return val.String(), nil
    case reflect.Bool:
    return strconv.FormatBool(val.Bool()), nil
    default: 
    return "", errors.New("Unsupported type")
    }
}


func (me Record) PutValue(key string, value reflect.Value) {
    switch (value.Kind()) {
        case reflect.Int, reflect.Int32, reflect.Int64:
            me.Putf(key, "%d", value.Int())
        case reflect.Uint, reflect.Uint32, reflect.Uint64:
            me.Putf(key, "%d", value.Uint())
        case reflect.Float32, reflect.Float64:
            me.Putf(key, "%f", value.Float())
        case reflect.String:
            me.Putf(key, "%s", value.String())
        default:
            me.Put(key, "???")
    }
}

func (me Record) PutStruct(prefix string, structure interface {}) {
    st := reflect.TypeOf(structure)
    vt := reflect.ValueOf(structure)
    monolog.Info("PutStruct: type %v value %v\n", st, vt)
    
    for i:= 0 ; i < st.NumField() ; i++ {
        field := st.Field(i)
        key := strings.ToLower(field.Name)
        value :=  vt.Field(i).String()
        me.Put(prefix + key, value)
    }
}


type Error struct {
    error   string
    lineno  int
}

func (me Error) Error() string {
    return fmt.Sprintf("%d: %s", me.Lineno, me.error)
}

func (me Error) Lineno() int {
    return me.lineno
}


type ParserState int

const (
    PARSER_STATE_INIT   ParserState = iota
    PARSER_STATE_KEY
    PARSER_STATE_VALUE
)

type RecordList []Record


func ParseReader(read io.Reader) (RecordList, error) {
    var records     RecordList
    var record      Record = make(Record)
    var err         Error
    lineno      := 0
    scanner     := bufio.NewScanner(read)
    var key     bytes.Buffer
    var value   bytes.Buffer
    
    
    for scanner.Scan() {
        lineno++
        line := scanner.Text()
        // End of record?
        if (len(line) < 1) || line[0] == '-' {
            // save the record and make a new one
            records = append(records, record) 
            record  = make(Record)
            // comment?
        } else if line[0] == '#' {
            continue; 
            // continue value?
        } else if line[0] == '\t' || line[0] == ' '|| line[0] == '+' {
            
            /* Add a newline unless + is used */
            if (line[0] != '+') {
                value.WriteRune('\n')          
            }
            
            // continue the value, skipping the first character
            value.WriteString(line[1:])            
            // new key
        } else if strings.ContainsRune(line, ':') {
            // save the previous key/value pair if needed
            if len(key.String()) > 0 {
                record[key.String()] = value.String()
            }
            
            key.Reset()
            value.Reset()

            parts := strings.SplitN(line, ":", 2)
                            
            key.WriteString(parts[0])
            if len(parts) > 1 {
               value.WriteString(parts[1])   
            }            
        // Not a key. Be lenient and assume this is a continued value.
        } else {
            value.WriteString(line)   
        }
    }
    
    // Append last record if needed. 
    if len(key.String()) > 0 {
        record[key.String()] = value.String()
    }
    
    if (len(record) > 0) {
        records = append(records, record)
    }

    

    if serr := scanner.Err(); serr != nil {
       err.lineno = lineno
       err.error  = serr.Error() 
       return records, err
    }
    
    return records, nil
    
}

func ParseFilename(filename string) (RecordList, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()
    return ParseReader(file)
}

func WriteField(writer io.Writer, key string, value string) {
    replacer := strings.NewReplacer("\n", "\n\t")    
    writer.Write([]byte(key))
    writer.Write([]byte{':'})    
    writer.Write([]byte(replacer.Replace(value)))
    writer.Write([]byte{'\n'})
}

func WriteRecord(writer io.Writer, record Record) {
    for key, value := range record {
        WriteField(writer, key, value);
    }
    writer.Write([]byte{'-', '-', '-', '-', '\n'})
}

func WriteRecordList(writer io.Writer, records RecordList) {
    for _, record := range records {
        WriteRecord(writer, record);
    }
}


func SaveRecord(filename string, record Record) (error) {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()
    WriteRecord(file, record)
    return nil
}



func SaveRecordList(filename string, records RecordList) (error) {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()
    WriteRecordList(file, records)
    return nil
}


/*
func ParseFile(file)
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
      // new record
      if line[0,2] == '--' 
        // Store last key used if any.
        if key          
          record[key] = value.chomp
          key = nil
        end  
        results << record
        record = {}
      elsif line[0] == '//'
      // Comments start with //
      elsif line[0] == ':'
      // a key/value pair
      key, value  = line[1,line.size].split(':', 2)
      record[key] = value.chomp
      key = value = nil
      elsif line[0, 2] == '..'
      // end of multiline value 
      record[key] = value.chomp
      key = value = nil
      elsif (line[0] == '.') && key.nil?
      // Multiline key/value starts here (but is ignored 
      // until .. is encountered)
      key   = line[1, line.size]
      key.chomp!
      value = ""
      // multiline value
      elsif key
          if line[0] == '\\'
            // remove any escapes
            line.slice!(0)
          end
          // continue the value
          value << line
      else
          // Not in a key, sntax error.
          errors << "//{lineno}: Don't know how to process line"
      end      
    end
    // Store last key used if any.
    if key      
      record[key] = value.chomp
    end  
    // store last record 
    results << record unless record.empty?
    return results, errors
  end  
  
  func load_filename(filename)
    results, errors = nil, nil, nil;
    file = File.open(filename, 'rt') rescue nil
    return nil, ["Could not open //{filename}"] unless file
    begin 
      results, errors = parse_file(file)
    ensure
      file.close
    end
    return results, errors
  end
  
  // Loads a Sitef fileas obejcts. Uses the ruby_klass atribute to load the object
  // If that is missing, uses defklass
  func load_objects(filename, defklass=nil)
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
  
  
  // Saves a single field to a file in Sitef format.
  func save_field(file, key, value)
    if value.is_a? String
      sval = value.dup
    else
      sval = value.to_s
    end
    if sval["\n"]
      file.puts(".//{key}\n")
      // Escape everything that could be misinterpreted with a \\
      sval.gsub!(/\n([\.\-\:\//\\]+)/, "\n\\\\\\1")
      sval.gsub!(/\A([\.\-\:\//\\]+)/, "\\\\\\1")
      file.printf("%s", sval)
      file.printf("\n..\n")
    else
      file.printf("://{key}://{sval}\n")
    end
  end
  
  func save_object(file, object, *fields)
    save_field(file, :ruby_class, object.class.to_s)
    fields.each do | field |
      value = object.send(field.to_sym)
      save_field(file, field, value)
    end
  end
  
  func save_record(file, record, *fields)
    record.each do | key, value |
      next if fields && !fields.empty? && !fields.member?(key)
      save_field(file, key, value)
    end
  end

  func save_file(file, records, *fields)
    records.each do | record |
      if record.is_a? Hash
        save_record(file, record, *fields)
      else 
        save_object(file, record, *fields)
      end
      file.puts("--\n")
    end
  end
  
  func save_filename(filename, records, *fields)
    results , errors = nil, nil
    file = File.open(filename, 'wt')
    return false, ["Could not open //{filename}"] unless file
    begin 
      save_file(file, records, *fields)
    ensure
      file.close
    end
    return true, []
  end
  
end

*/

