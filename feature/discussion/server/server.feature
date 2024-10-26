Feature: Server

Scenario: user authorization
    Given there is a user Dima
    When Dima authorises
    Then he can do it

  Scenario: new meeting
    Given there is a user Dima
    When Dima creates a new meeting
    Then he can do it

  Scenario: list meetings
    Given there is a user Dima
    And he has meetings A and B
    And there is a user Alex
    And he has meetings C and D
    When Dima queries meetings
    Then he sees his meetings A and B
    But he does not see meetings C and D, which belong to Alex
