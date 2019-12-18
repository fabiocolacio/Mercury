defmodule Mercury.Auth.Token do
  use Guardian, otp_app: :mercury

  def subject_for_token(uid, _claims) do
    {:ok, uid}
  end

  def resource_from_claims(%{"sub" => uid}) do
    {:ok, uid}
  end
end
