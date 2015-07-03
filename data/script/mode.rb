
# different modes of the client state
module Mode
  def do_command(str)
    met = "on_#{@state}".to_sym
    if self.respond_to?(met) 
      self.send(met, str)
    else
      log "Unknown state #{@state}"
      @state = :login
    end
  end
end

