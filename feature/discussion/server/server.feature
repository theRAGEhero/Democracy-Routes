Feature: Server

    Scenario: authorization
        Given there is a user Dima
        When Dima authorises
        Then he can do it

    Scenario: meeting
        Given there is a user Dima
        When Dima creates a new meeting
        Then he can do it
