defmodule Mercury.Router.Protected do
  use Plug.Router

  plug Mercury.Auth.Pipeline
  plug :match
  plug :dispatch

  get "/recv" do
    send_resp(conn, 200, "/recv")
  end

  post "/send" do
    send_resp(conn, 200, "/send")
  end

  match _ do
    send_resp(conn, 404, "Bad endpoint.")
  end
end
