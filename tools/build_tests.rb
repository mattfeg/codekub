#!/usr/bin/env ruby
# Builds tests/NN.in and tests/NN.out for exercises.
#
# Inputs come from the exercise's generator.rb, which prints a JSON array of
# input strings (no generator means a single empty input). Expected outputs come
# from running reference.rb on each input, so no expected output is ever typed
# by hand. The first `samples` outputs must appear verbatim in statement.md.
#
# Usage: tools/build_tests.rb [exercise_dir ...]   (default: every exercise)
require "json"
require "open3"
require "fileutils"

ROOT = File.expand_path("..", __dir__)

def run_ruby(script, dir, stdin = "")
  out, err, status = Open3.capture3("ruby", script, stdin_data: stdin, chdir: dir)
  abort "#{File.join(dir, script)} failed:\n#{err}" unless status.success?
  out
end

def build(dir)
  meta = JSON.parse(File.read(File.join(dir, "meta.json")))
  inputs = File.exist?(File.join(dir, "generator.rb")) ? JSON.parse(run_ruby("generator.rb", dir)) : [""]
  abort "#{dir}: generator produced no inputs" if inputs.empty?

  tests = File.join(dir, "tests")
  FileUtils.rm_rf(tests)
  FileUtils.mkdir_p(tests)
  width = [2, inputs.size.to_s.size].max
  outputs = inputs.each_with_index.map do |input, i|
    name = (i + 1).to_s.rjust(width, "0")
    output = run_ruby("reference.rb", dir, input)
    File.write(File.join(tests, "#{name}.in"), input)
    File.write(File.join(tests, "#{name}.out"), output)
    output
  end

  statement = File.read(File.join(dir, "statement.md"))
  outputs.first(meta.fetch("samples", 1)).each_with_index do |output, i|
    abort "#{dir}: sample #{i + 1} output is not in statement.md" unless statement.include?(output.strip)
  end

  # Imported problems ship their official tests: the reference must agree with
  # every official output, which also proves the reference is correct.
  official_path = File.join(dir, "official_tests.json")
  if File.exist?(official_path)
    expected = JSON.parse(File.read(official_path)).map { |t| [t["input"].split, t["output"].split] }.to_h
    inputs.zip(outputs).each_with_index do |(input, output), i|
      want = expected[input.split]
      abort "#{dir}: test #{i + 1} disagrees with the official output" if want && want != output.split
    end
  end
  puts format("  %-32s %2d teste(s)", File.basename(dir), inputs.size)
end

# By default, every exercise that has a reference solution (imported beecrowd
# problems without one keep the tests they came with).
dirs = ARGV.empty? ? Dir[File.join(ROOT, "track/*/*/reference.rb")].map { |m| File.dirname(m) } : ARGV
dirs.sort.each { |d| build(File.expand_path(d)) }
