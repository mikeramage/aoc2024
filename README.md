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

## Day 5
Go, why do you have no set in your standard collections? Yes, I know I can just use a map, but sometimes I want a set. Is it too much to ask? It's mildly irritating but also quite fun getting used to Go's quirks. As I mentioned above, the syntax feels to me like a hybrid of C-ish languages that I know (C, C++, Java, Rust) and Python, and I sometimes assume Go will behave like one or the other, before smacking me in the chops because Go is Go and it's a bit different from any of these. 

It's a bit like that time I annoyed a bunch of Italians. I already spoke Spanish and a some French, so I spent a couple of weeks learning some basic Italian, and then went on holiday to Italy, possessed of the hideously arrogant misconception that I was fluent in the language. You know, if you don't know the word it's bound to be some Italianisation of whatever I'd say in Spanish/French, right? Am I right??? So yeah, I just swanned around Rome, happily spamming the locals with my bastard pidgin until being humbled by a bunch of confused looks and the occasional rude gesticulation and accompanying stream of fruity invective. AoC is definitely good for learning the idioms of a new coding language (particularly if that language happens to be Rust: I'm looking at you, Madam Compiler) in the same way an that an angry Roman flipping his fingers off his chin and telling you where to go (inside your own bottom, apparently) as if you've just insulted his entire family is quite motivating for _actually bothering to get off your lazy arse and learn_ the true idioms of his spoken language. (I should probably say I'm exaggerating (er ... _lying_) for comic effect. In reality, all the Italians I met on that holiday were really friendly and usually found my "Italian" hilarious - if occasionally mystifying. At least I appeared to be _trying_ to speak it instead of shouting English at them slowly.)

Anyhow, Day 5 was straightforward enough. A basic dependency graph problem. I just used a map of pages to their set (I mean - grrrr! - _map_) of prerequisites and it fell out reasonably quickly (once I'd fixed the sloppy bugs I introduced along the way of course). 

## Day 6
Brute forced this one.

## Day 7
Brute forced this one too. 

Yeah.

## Day8 
Bunch of maps. Counting antinodes. Little picture on paper of adding/subtracting vectors so I don't get too befuddled. Fine. Whatever Can't really be bothered looking for more optimal solutions or making my code nice or anything. Got the answer. It's a Sunday. Happy. Yeah ... 

... Well, a lot happier than I am about the water pipe leaking in my bathroom. The pipe is copper and some imcomprehensible clown of a contractor (or perhaps it was some recklessly over-ambitious previous home owner somewhere low down on the DIY Dunning-Kruger curve; there are some now-disabled "interesting" DIY electrics in the attic as well as some highly dubious carpentry in the kitchen so I wouldn't be surprised) went and partially embedded it in the concrete underflooring when they were laying the lovely faux wood floor panels (who the hell puts a "wooden" floor in their bathroom anyway. It's a bathroom. Use tiles you idiots!) Of course, the concrete seems to have reacted with the copper over the years and has corroded it, resulting in a burst pipe, a puddle on my floor, who knows what behind the wooden panelling, and a day of stress on hold to the insurance company who are currently overwhelmed by the volume of claims arising from Storm Darragh. Of course they need to cut the wall away and possibly get a massive chisel out to hack at the concrete and it's costing a minor fortune. I'm paying one way or another. Either I pay for it all, or the insurance company does and I pay the excess and an increased premium the rest of my life. Anyway, rant over. Might actually talk about something actually relevant to advent of code tomorrow. Yeah, right! 
