Feature: Server

    Scenario: authorization
        Given there is a user Dima
        When Dima authorises
        Then he can do it
