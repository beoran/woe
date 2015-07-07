require 'atto'
include Atto::Test

require_relative '../../lib/woe/settings' 

assert { Woe::Settings }

assert { Woe::Settings.port         == 7000           }
assert { Woe::Settings.data_dir     == 'data'         }
assert { Woe::Settings.var_dir      == 'data/var'     }
assert { Woe::Settings.script_dir   == 'data/script'  }

ARGV << '-p'
ARGV << '7777'

ARGV << '-d'
ARGV << '/var/woe/data'

p ARGV

assert { Woe::Settings.parse_args }


assert { Woe::Settings.port         == 7777           }
assert { Woe::Settings.data_dir     == '/var/woe/data'         }
assert { Woe::Settings.var_dir      == '/var/woe/data/var'     }
assert { Woe::Settings.script_dir   == '/var/woe/data/script'  }


# Undefining class methods
# ec = class << Kernel ; self; end
# (Kernel.methods - BasicObject.methods).each do |m| ; ec.class_eval { remove_method m } ; end
# Removing constants
# Object.instance_eval { remove_const :ARGF }

