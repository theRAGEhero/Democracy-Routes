Feature: Server

  Scenario: new user
    Given a new user Dima is added
    Then user Dima exists

  Scenario: user authorization
    Given there is a user Dima
    When Dima authorises
    Then he can do it

  Scenario: new meeting
    Given there is a user Dima
    When Dima creates a new meeting
    Then he can do it
