defmodule Mercury.Auth.Pipeline do
  use Guardian.Plug.Pipeline,
    otp_app: :mercury,
    module: Mercury.Auth.Token,
    error_handler: Mercury.Auth.ErrorHandler

  @claims %{iss: Application.get_env(:mercury, :issuer)}

  plug Guardian.Plug.VerifySession, claims: @claims
  plug Guardian.Plug.VerifyHeader, claims: @claims
  plug Guardian.Plug.EnsureAuthenticated
end
