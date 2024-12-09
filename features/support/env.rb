Before do
  # we don't need to send the integrity header
  Maze.config.enforce_bugsnag_integrity = false
  $address = nil
  steps %(
    When I configure the maze endpoint
  )
end
