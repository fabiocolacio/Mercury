defmodule Mercury.Repo do
  use Ecto.Repo,
    otp_app: :mercury,
    adapter: Ecto.Adapters.MyXQL
end
