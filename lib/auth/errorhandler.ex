defmodule Mercury.Auth.ErrorHandler do
  @behaviour Guardian.Plug.ErrorHandler

  require Logger

  @impl Guardian.Plug.ErrorHandler
  def auth_error(_conn, {_type, _reason}, _opts) do
    Logger.info("Authentication error")
  end
end
