defmodule Mercury.Router do
  use Plug.Router

  plug :match
  plug :dispatch

  post "/register" do
    send_resp(conn, 200, "/register")
  end

  get "/login" do
    send_resp(conn, 200, "/login (get)")
  end

  post "/login" do
    send_resp(conn, 200, "/login (post)")
  end

  forward "/user", to: Mercury.Router.Protected

  match _ do
    send_resp(conn, 404, "Bad endpoint.")
  end
end