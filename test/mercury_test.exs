defmodule MercuryTest do
  use ExUnit.Case
  doctest Mercury

  test "greets the world" do
    assert Mercury.hello() == :world
  end
end
