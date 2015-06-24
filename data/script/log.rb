
# log function, with formatting support and auto newline.
def log(fmt, *args)
  if fmt && !(args.empty?) 
    woe_log(fmt.format(*args)) 
  else
    woe_log(fmt.to_s) 
  end
end

log("log.rb script loaded")

