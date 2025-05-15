# bet_evaluator

This project is a simple command-line application that evaluates betting outcomes for volleyball and cricket matches based on provided pre-match odds and actual results.

## How to Run

1.  **Prerequisites:**

    - Go (version 1.18 or later for generics support) installed on your system.
    - Git

2.  **Clone the Repository:**

    ```bash
    git clone https://github.com/Kirubel-Enyew27/bet_evaluator

    cd bet_evaluator
    ```

3.  **Create Data Folders and Files:**

    - Create two subdirectories within the project root: `volleyball` and `cricket`.
    - Place your pre-match and result data JSON files in the respective folders if you want to change/modify the given data:
      - `volleyball/volleyball_prematch.json`
      - `cricket/cricket_prematch.json`
      - `volleyball/volleyball_result.json`
      - `cricket/cricket_result.json`
        **Note:** The structure of these JSON files should match the `models` package definitions used in the Go code.

4.  **Run the Application:**

    ```bash
    go run main.go
    ```

5.  **Follow Prompts:**
    The application will prompt you to enter the sport you want to evaluate (`volleyball` or `cricket`) or `exit` to quit. Enter your choice and press Enter.

    ```
    Enter your choice from sport betting [volleyball, cricket] or 'exit' to quit: volleyball
    ```

    or

    ```
    Enter your choice from sport betting [volleyball, cricket] or 'exit' to quit: cricket
    ```

    The application will then load the corresponding pre-match and result data and print the evaluation of different betting markets.

6.  **Exiting:**
    To exit the application, enter `exit` at the prompt.

## Research on Betting Markets

This section provides a brief overview of how each betting market works, as implemented in the code.

### Volleyball Betting Markets

1.  **Match Winner (1X2):**

    - **Description:** Betting on the outright winner of the volleyball match.
    - **How it Works:** The code extracts odds for the home team (Header "1") and the away team (Header "2"). It then compares the final set score (`SS`) from the result data to determine the winner. If the home team's set score is higher, a bet on the home team wins, and vice versa. Draws are not possible in standard volleyball.

2.  **Correct Set Score:**

    - **Description:** Betting on the exact final set score of the match (e.g., 3-0, 3-1, 2-3).
    - **How it Works:** The code iterates through all the offered correct set score odds in the pre-match data. It compares the predicted score (extracted from `odd.Name` and `odd.Header` indicating home or away win) with the actual final set score (`result.SS`). If they match, the bet is considered "WON".

3.  **Total Points:**

    - **Description:** Betting on the total number of points scored by both teams combined in the entire match being over or under a specified line.
    - **How it Works:** The code calculates the total points by summing the points scored by both teams in each set from the result data (`result.Scores`). It then compares this total with the "Over" (`O`) or "Under" (`U`) line specified in the pre-match odds (`odd.Handicap`).

4.  **Handicap:**

    - **Description:** Betting on the outcome of the match after one team is given a virtual advantage (handicap).
    - **How it Works:** The code calculates the set difference between the home and away teams (`homeScore - awayScore`). It then evaluates handicap bets based on whether the handicap is positive (applied to the away team) or negative (applied to the home team).

5.  **Double Chance:**
    - **Description:** Betting on two of the three possible outcomes (Home Win or Draw, Away Win or Draw, Home Win or Away Win). Note that a draw is not possible in volleyball, so only Home Win or Away Win (12) is evaluated here.
    - **How it Works:** The code checks if the final set score resulted in a win for either the home or away team (i.e., the scores are not equal). If they are not equal, a bet on "Home or Away" is considered "WON".

### Cricket Betting Markets

1.  **Match Winner:**

    - **Description:** Betting on the outright winner of the cricket match.
    - **How it Works:** The code extracts odds for team "1" (assumed to be home) and team "2" (assumed to be away). It compares the final score (parsed from `result.SS`) to determine the winner.

2.  **Most Match Sixes:**

    - **Description:** Betting on which team will hit the most sixes in the match.
    - **How it Works:** The code makes a simplified assumption that the team scoring more runs hit more sixes. If the home team's runs are higher, a bet on "1" (home) wins, and similarly for the away team ("2"). A tie in runs results in a "Tie" bet winning. **Note:** This is a simplification as run count doesn't directly correlate to sixes. A more accurate implementation would require specific data on the number of sixes hit by each team.

3.  **Most Match Fours:**
    - **Description:** Betting on which team will hit the most fours in the match.
    - **How it Works:** Similar to the "Most Match Sixes" market, the code uses the final run count as a proxy for the number of fours hit.

**Note:** The cricket betting evaluation in this project makes some simplifying assumptions, particularly for the "Most Sixes" and "Most Fours" markets, due to the lack of direct sixes and fours data in the provided structure.
