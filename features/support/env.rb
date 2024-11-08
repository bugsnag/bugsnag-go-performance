Before do
  # we don't need to send the integrity header
  Maze.config.enforce_bugsnag_integrity = false
  # TODO - can be restored after https://smartbear.atlassian.net/browse/PIPE-7498
  Maze.config.skip_default_validation('trace')
  $address = nil
  steps %(
    When I configure the maze endpoint
  )
end
