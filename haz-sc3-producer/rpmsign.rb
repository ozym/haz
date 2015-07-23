#!/usr/bin/env ruby

require 'pty'
require 'expect'

key=ARGV[0]
passphrase=ARGV[1]
rpms=ARGV[2..-1]

command = "rpm --define '_signature gpg' --define '_gpg_name #{key}' --addsign " + rpms.join(' ')
puts command

output = ''

PTY.spawn(command) do |r,w,p|

  w.sync = true
  $expect_verbose = true

  r.expect(/Enter pass phrase: /)
  sleep(0.5)
  w.puts(passphrase)
  begin
    r.each { |l| output += l }
  rescue Errno::EIO
  end
end

puts output

exit
