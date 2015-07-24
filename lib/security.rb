#
# Woe security related helper functions. 
#


CRYPT_MAKE_SALT_AID = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789./"

# Generates salt for use by crypt.
def crypt_make_salt  
  c1 = CRYPT_MAKE_SALT_AID[rand(CRYPT_MAKE_SALT_AID.length)] 
  c2 = CRYPT_MAKE_SALT_AID[rand(CRYPT_MAKE_SALT_AID.length)] 
  return c1 + c2
end

# Crypt with salt generation.
def crypt(pass, salt = nil) 
  salt = crypt_make_salt unless salt
  return pass.to_s.crypt(salt)
end

# Challenge crypt password trypass against the hash hash
def crypt_challenge?(trypass, hash) 
  salt = hash[0, 2]
  tryhash = trypass.to_s.crypt(salt)
  return tryhash == hash
end
