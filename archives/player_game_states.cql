CREATE TYPE chesshouzs.skill_position (
    row INT, 
    col INT
);

CREATE TYPE chesshouzs.skill_status (
    position FROZEN<skill_position>, 
    duration_left INT
);

CREATE TYPE chesshouzs.state (
    duration_left INT, 
    list LIST<FROZEN<skill_status>>
);

-- Create the table
CREATE TABLE chesshouzs.player_game_states (
    player_id TEXT,
    game_id TEXT,
    buff_state MAP<TEXT, FROZEN<state>>,
    debuff_state MAP<TEXT, FROZEN<state>>,
    PRIMARY KEY ((player_id, game_id))
);