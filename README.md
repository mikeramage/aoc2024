# Advent of Code 2024

Well, here we Go!

... Sorry. I'll get my coat.

## Day 1
Another year, another set of puzzles, another ludicrous README, a new programming language. 

Typical Day 1 puzzle, easing you in oh so gently, the calm before the storm. Spent most of my 
time setting up a framework for running the code. Go is a curious beast. Considerably less 
painful than Rust to learn, yet missing a lot of its pretty features. Syntax somewhere between C 
and Python. Error handling seems very verbose (i.e. if err != nil everywhere). For this puzzle I wanted to 
zip() a couple of lists together a la Python, but it doesn't exist, so I've got a C-style for loop 
over indices, which feels a bit old-skool, but whatever. Reading lines from a file was a bit funky.
Open a file, handle failure, create a scanner, scan the lines. I dunno, it all feels a bit like teaching 
my parents how the internet works. A lot of words, seemingly unnecessary repetition, not really sure if 
they (or I) have really understood it all. 

Still, 2 stars is 2 stars and I quite enjoyed mucking about with cobra-cli. Onwards and upwards. Roll on 
Day 2! 

## Day 2
Dear oh dear oh dear. Weak. This was extremely painful to watch for the omnipresent, obnoxiously critical, entirely imaginary observer leering over my shoulder (imagine some dusty old Professor Snape figure sneering away). I finished ages ago and he's still tutting away, whispering snide comments and disparaging remarks in my ear. Perhaps staying up to 2am last night wasn't the brightest of ideas, or perhaps my mind is simply dimming as the cobwebs encroach, marking my steady progress into the depths of middle age. No! I'm not dead yet, dammit! I CAN STILL DO THIS!

Anyway, I got to the answer eventually by way of a few coffees and some confused meanderings around the misty peripheries of reason and logic, haphazardly stumbling onto a successful solution like one of those infinitely many monkeys with a typewriter. Turns out I was trying to be a bit too "clever" (admittedly, based on how this puzzle went, I'm taking quite the liberties with the word), going for some Leetcode-style highly-optimised single-pass solution and getting bewildered by edge cases. It's not a job interview you doofus. Just use brute force like you always do! Ugh. 

After solving it myself, I had a look on YouTube at Jonathan Paulson's Python solution, which took a different approach. After I stopped crying (he solved it in 3 minutes 55 seconds despite having some annoying connectivity issues with the server that probably cost him around a minute, whereas it took me ... well, never you mind how long it took me. It's not a competition. Jeez, just leave it, will you?), I re-implemented his algorithm in Go (the "Alt" functions), which was clunky but educational - got to look at some iter and slices package stuff.

Tomorrow is another day.

## Day 3
Still not exactly showering myself in glory. Easier than yesterday, but instead of using a simple regex for part 2, I went down a rabbit hole of finding start and stop indices manually and then slicing and dicing on those instead of letting the regex engine doing it for me. But I got the answer both ways and used more Go features than I otherwise would. 

Regex in Go feels a bit incomplete, particularly the handling of named capturing groups (unless I'm missing something, you need to map the names to the result values manually to be able to refer to them by name, which kind of defeats the point). 

## Day 4
Out all day marshalling school kids around public transport and a busy museum (and yes, I was _supposed_ to be marshalling the kids. It was sanctioned by the school. I didn't just decide to round up some random children and move them around the city all day for something to do). And so I didn't get started till after dinner. And I'm exhausted since herding cats is tiring. So it didn't go that well, OK? 

Anyway, I had some idea of organising the input into a 2D character array, creating 1D arrays of string "views" at the different angles (diagonal is quite fun), and doing a string count of XMAS and SAMX for each string in the views for the 4 different orientations. I'm quite glad I started down this track because I came across "runes" in the process, which is a cool name for unicode codepoints. But I got fed up because it turned out to be a bit of a faff and, like I said, I'm tired. And more importantly, I sensed it was unlikely that that approach would work for part 2. How right I was. So I came up with a very dull, workmanlike implementation that might as well have been written in C. Bo-ring!