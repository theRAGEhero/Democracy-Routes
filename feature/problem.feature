Feature: Problem solving
    A person can solve a problem by asking experts.

    Scenario: a new problem creation
        Given there is a person named "Alex"
        When "Alex" creates a problem "How to build a decision making system"
        And there is no similar problems yet
        Then the problem is created
        And the new "Route" is added to the problem
        And the new "Round Table" is added to the route

    Scenario Outline: ability to reuse a route of an already existing problem
        Given there is a problem "How to build a decision making system"
        And this problem has a "Route"
        And this "Route" has a "Round Table"
        When "Dima" creates a similar problem "How to build an automated decision making system"
        Then he receives a suggestion to reuse the "Route" for the existing problem
        And he receives a suggestion to reuse the "Round Table" for the existing route
        And he makes a <decision>
        And he gets the <result>
        Examples:
            | decision | result      |
            | accept   | reuse_route |
            | decline  | new_problem |

    Scenario Outline: ability to create a new round table
        Given there is a problem "How to build a decision making system"
        And this problem has round tables "Round Table-1" and "Round Table-2"
        When "Dima" creates a new round table for this problem
        Then he can choose to descend it from the <round table>
        And he creates "Round Table-3" descended from the <round table>
        Examples:
            | round table |
            | 1           |
            | 2           |

    Scenario: grouping similar round tables to a virtual route
        Given there is a problem "How to build a decision making system"
        And there is a "Route 1" for this problem
        And there is a "Round Table 1" for this route
        And there is a "Route 2" for this problem
        And there is a "Round Table 2" for this route
        When "Round Table 1" and "Round Table 2" have similar topics
        Then they becomes grouped in to the same "Virtual Route"

    Scenario: Searching experts
        Given a person named "Alex"
        And a problem "How to build decision making system"
        And experts in this problem "Dima" and "Carlo"
        And experts in another problem "Shaun" and "Bob"
        When "Alex" searches for experts in "How to build decision making system"
        Then he finds "Dima" and "Carlo"
        But he does not find "Shaun" and "Bob"

    Scenario: Problem discussion
        Given a person named "Alex"
        And a problem "How to build a decision making system"
        And experts in this problem "Dima", "Carlo"
        When "Alex" starts discussion of the problem
        Then "Dima", "Carlo" can discuss it using "video rooms"
        And the "automatic transcription" of this discussion is saved
        And the "automatic summary" of this discussion is saved
        And the problem discussion becomes finished

    Scenario: Correcting the summary of a finished problem discussion
        Given there is a finished problem discussion
        When the patricipants want to change it's summary
        Then they can do it

    Scenario: Signing a finished problem discussion
        Given there is a finished problem discussion
        When the patricipants sign it
        Then the discussion becomes signed
        And the patricipants can no longer change it's summary
