
# Model and handle the clients of the server
class Timer 
  
  
  attr_reader :id
  attr_reader :name
  
  def initialize(name, &block)
    @name     = name
    @id       = nil
    @block    = block
  end
  
  def set(value = 1.0, interval = 1.0)
     Woe::Server.set_timer(@id, value, interval);
  end
  
  def start(value = 1.0, interval = 1.0)
    @id = Woe::Server.new_timer()
    return @id if (@id < 0)
    self.set(value, interval);
    return @id
  end
  
  def self.add(name, value, interval, &block)
    timer = Timer.new(name, &block)
    return nil if timer.start(value, interval) < 0 
    @timers ||= {}
    @timers[timer.id] = timer
  end
  
  def self.get(id)
    @timers ||= {}
    @timers[id]  
  end
  
  def self.get_by_name(name) 
    @timers ||= {}
    @timers.select { |t| t.name == name }.first 
  end
  
  def self.remove(id)    
    @timers[id] = nil
  end
  
  def on_timer()
    @block.call
  end
  
  def self.on_timer(id) 
    timer = self.get(id)
    timer.on_timer if timer
  end
  
end

log "Mruby timer script loaded OK."

