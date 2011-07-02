#!/usr/bin/ruby

def update_urls(s, tag)
	s.gsub(/\/(tree|blob)\/master\//, "/\\1/#{tag}/")
end

def sys(*command)
	def add_quotes(ss)
		ss.map { |s| (s.include? ' ') ? '"'+s+'"' : s }.join(' ')
	end
	system(*command) or raise "command failed: #{add_quotes(command)}"
end

version = ARGV[0]
raise "version number must be in format 'x.y.z' but was '#{version}'" unless version =~ /\A\d+\.\d+\.\d+\z/m
README_FILE = 'README.md'

prev_readme = IO.read(README_FILE)
readme_parts = prev_readme.partition(/\*\*\d+\.x\.x \(20..-xx-xx\)\*\*/)

# release readme
app = File.basename(Dir.getwd)
tag = "#{app}-#{version}"
release_date = Time.now.strftime('%Y-%m-%d')
release_title = "**#{version} (#{release_date})**"
release_readme = update_urls(readme_parts[0], tag) + release_title + update_urls(readme_parts[2], tag)

# next readme
major = version.split('.')[0]
year = Time.now.strftime('%Y')
next_title = "**#{major}.x.x (#{year}-xx-xx)**\n\n- ...\n\n"
next_readme = readme_parts[0] + next_title + release_title + readme_parts[2]

# commit and tag release
raise "must be in master branch" unless `git branch` =~ /^\* master$/
File.open(README_FILE, 'w') { |f| f.write(release_readme) }
sys('git', 'add', README_FILE)
sys('git', 'commit', '-m', "Release #{version}")
sys('git', 'tag', tag)

# update release branch
sys('git', 'checkout', 'release')
sys('git', 'merge', '--ff-only', tag)
sys('git', 'checkout', 'master')

# commit next iteration
File.open(README_FILE, 'w') { |f| f.write(next_readme) }
sys('git', 'add', README_FILE)
sys('git', 'commit', '-m', "Prepare for next development iteration")

# manual check before pushing the release
sys('gitk', '--all')
puts "\nProceed with release? (yes/no)"
if STDIN.gets.chomp == 'yes'

	# push
	sys('git', 'push', 'origin', 'master:master')
	sys('git', 'push', 'origin', 'release:release')
	sys('git', 'push', '--tags', 'origin')

	puts "\nRelease done."
else

	# undo tag
	sys('git', 'tag', '-d', tag)
	
	# undo release branch
	sys('git', 'checkout', 'release')
	sys('git', 'reset', '--hard', 'origin/release')
	
	# undo master branch
	sys('git', 'checkout', 'master')
	sys('git', 'reset', '--hard', 'HEAD~2')

	puts "\nRelease aborted."
end
