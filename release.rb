#!/usr/bin/ruby

def update_urls(s, tag)
	s.gsub(/\/(tree|blob)\/master\//, "/\\1/#{tag}/")
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
system('git', 'add', README_FILE) or raise
system('git', 'commit', '-m', "Release #{version}") or raise
system('git', 'tag', tag) or raise

# update release branch
system('git', 'checkout', 'release') or raise
system('git', 'merge', '--ff-only', tag) or raise
system('git', 'checkout', 'master') or raise

# commit next iteration
File.open(README_FILE, 'w') { |f| f.write(next_readme) }
system('git', 'add', README_FILE) or raise
system('git', 'commit', '-m', "Prepare for next development iteration") or raise

# manual check before pushing the release
system('gitk', '--all') or raise
puts "\nProceed with release? (yes/no)"
if STDIN.gets.chomp == 'yes'

	# push
	system('git', 'push', 'origin', 'master:master') or raise
	system('git', 'push', 'origin', 'release:release') or raise
	system('git', 'push', '--tags', 'origin') or raise

	puts "\nRelease done."
else

	# undo tag
	system('git', 'tag', '-d', tag) or raise
	
	# undo release branch
	system('git', 'checkout', 'release') or raise
	system('git', 'reset', '--hard', 'origin/release') or raise
	
	# undo master branch
	system('git', 'checkout', 'master') or raise
	system('git', 'reset', '--hard', 'HEAD~2') or raise

	puts "\nRelease aborted."
end

