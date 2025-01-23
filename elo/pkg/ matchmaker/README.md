# matchmaker algorithm

## Terms
- Ticket: A ticket is both the unit and the subject of matchmaking. A player or team (group of players) associated with related data (preferred maps etc) that wait to be matched is included.
- Queue: collection of tickets to be matched. Set of rules that controls how tickets are matched.
- Rule: A constraint on which tickets are eligible to match. The matchmaking algorithm searches for set of tickets that satisfy all the rules defined by a queue to create a match.
- Attribute: data associated with a player
- Match: output of matchmaking process

## Goal
- Match tickets in pairs such that the total MMR diff across all pairs is minimized.

## Algorithm
> Weighted Greedy Pairing Algorithm

The scoring function can be represented as:


$$
\text{Score}(A, B) = \alpha \cdot |\text{MMR}_A - \text{MMR}_B| + \beta \cdot (\text{WaitingTime}_A + \text{WaitingTime}_B)
$$

Where:
- $\alpha$ is the weight for MMR difference.
- $\beta$ is the weight for waiting time.
- A and B are the two players being considered for a match.


The goal is to minimize the score for each match.

## Ref
[MS) Matchmaking](https://learn.microsoft.com/en-us/gaming/playfab/features/multiplayer/matchmaking/)