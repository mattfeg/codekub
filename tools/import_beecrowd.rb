#!/usr/bin/env ruby
# Creates the track's beecrowd problems from a clone of
# https://github.com/potigol/beecrowd.
#
# Each exercise gets the title, category and link from docs/problemas/*.md and
# its tests from the statement examples in src/*/exemplos/. The statement
# itself stays on beecrowd: codekub opens it in the browser (`o` in the panel).
# Problems without examples, and SQL ones, are skipped.
#
# Only missing phases and exercises are created: anything already in the track
# (a written statement, new tests) is never overwritten.
#
# Usage: tools/import_beecrowd.rb <potigol-beecrowd-clone> [track_dir]
#        (default track_dir: the repository's track/)
require "json"
require "fileutils"
require "set"

ROOT = File.expand_path("..", __dir__)

PHASES = [
  # [potigol category, phase slug, phase title (short: it fits the side panel), description]
  ["Iniciante", "iniciante", "Iniciante", "Entrada, saída, condicionais e laços"],
  ["Ad-Hoc", "adhoc", "Ad-Hoc", "Simulação, datas e problemas gerais"],
  ["Strings", "strings", "Strings", "Manipulação de texto"],
  ["Matemática", "matematica", "Matemática", "Sistemas numéricos, primos e números grandes"],
  ["Estruturas e Bibliotecas", "estruturas", "Estruturas", "Estruturas e bibliotecas: filas, pilhas, ordenação e mapas"],
  ["Paradigmas", "paradigmas", "Paradigmas", "Programação dinâmica, busca binária e gulosos"],
  ["Grafos", "grafos", "Grafos", "Busca, caminhos mínimos e árvores"],
  ["Geometria Computacional", "geometria", "Geometria", "Geometria computacional: pontos, linhas e polígonos"],
].freeze

def slugify(s)
  s.unicode_normalize(:nfd).gsub(/\p{Mn}/, "").downcase.gsub(/[^a-z0-9]+/, "-").gsub(/\A-|-\z/, "")
end

def read_text(path)
  File.read(path, encoding: "UTF-8").gsub("\r\n", "\n")
end

# docs/problemas/*.md: "## 1001 - Title", then the statement link and category.
def problems(clone)
  found = {}
  Dir[File.join(clone, "docs", "problemas", "*.md")].sort.each do |file|
    read_text(file).split(/^## (?=\d+ - )/)[1..-1].to_a.each do |section|
      head = section[/\A(\d+) - (.+)$/] && [$1, $2.strip]
      url = section[/^Descrição do Problema: \[[^\]]*\]\((\S+?)\)/, 1]
      category = section[/^Categoria: \[([^\]]+)\]/, 1]
      next unless head && url && category
      found[head[0]] ||= { id: head[0], title: head[1], url: url, category: category }
    end
  end
  found.values
end

def examples(clone, id)
  Dir[File.join(clone, "src", "*", "exemplos", "#{id}-*.in")]
    .sort_by { |f| f[/-(\d+)\.in\z/, 1].to_i }
    .map { |f| [f, f.sub(/\.in\z/, ".out")] }
    .select { |_, out| File.exist?(out) }
    .map { |i, o| [read_text(i), read_text(o)] }
end

def block(text)
  "```\n#{text.chomp}\n```\n"
end

def statement(p, samples)
  parts = ["# #{p[:title]}\n",
           "O enunciado está no beecrowd: aperte `o` para abri-lo no navegador.\n",
           "<#{p[:url]}>\n",
           "Os testes são só os exemplos do enunciado: passar neles não garante que a\nsolução esteja certa no beecrowd.\n",
           "## Exemplos\n"]
  samples.each do |i, o|
    parts << (i.strip.empty? ? "Saída:\n\n#{block(o)}" : "Entrada:\n\n#{block(i)}\nSaída:\n\n#{block(o)}")
  end
  parts.join("\n")
end

clone = ARGV[0] or abort "uso: tools/import_beecrowd.rb <clone-do-potigol-beecrowd> [pasta_da_trilha]"
abort "#{clone} não parece um clone de potigol/beecrowd" unless Dir.exist?(File.join(clone, "docs", "problemas"))
out = File.expand_path(ARGV[1] || File.join(ROOT, "track"))
existing = Dir[File.join(out, "*", "*", "meta.json")].map { |m| File.basename(File.dirname(m))[/\A\d+/] }.to_set

by_category = problems(clone).group_by { |p| p[:category] }
skipped = Hash.new(0)
total = 0
PHASES.each_with_index do |(category, slug, title, description), i|
  dir = File.join(out, format("%02d-%s", i + 1, slug))
  FileUtils.mkdir_p(dir)
  phase_file = File.join(dir, "phase.json")
  unless File.exist?(phase_file)
    File.write(phase_file, JSON.pretty_generate(
      id: format("%02d", i + 1), slug: slug, title: title, description: description,
      open: true, topics: [{ slug: slug, title: category }]
    ) + "\n")
  end
  by_category.fetch(category, []).sort_by { |p| p[:id].to_i }.each do |p|
    if existing.include?(p[:id])
      skipped["já existiam"] += 1
      next
    end
    samples = examples(clone, p[:id])
    if samples.empty?
      skipped["sem exemplos"] += 1
      next
    end
    ex_dir = File.join(dir, "#{p[:id]}-#{slugify(p[:title])}")
    FileUtils.mkdir_p(File.join(ex_dir, "tests"))
    File.write(File.join(ex_dir, "meta.json"), JSON.pretty_generate(
      id: p[:id], slug: slugify(p[:title]), title: p[:title], topic: slug, level: 0, concepts: [],
      time_limit: 1.0, samples: samples.size,
      source: { name: "beecrowd #{p[:id]}", url: p[:url], dataset: "potigol/beecrowd (exemplos)" }
    ) + "\n")
    File.write(File.join(ex_dir, "statement.md"), statement(p, samples))
    samples.each_with_index do |(i_text, o_text), k|
      File.write(File.join(ex_dir, "tests", format("%02d.in", k + 1)), i_text)
      File.write(File.join(ex_dir, "tests", format("%02d.out", k + 1)), o_text)
    end
    total += 1
  end
end
(by_category.keys - PHASES.map(&:first)).each { |c| skipped["categoria #{c}"] += by_category[c].size }

puts "#{total} problemas criados em #{out}"
skipped.each { |why, n| puts "  #{n} ignorados (#{why})" }
