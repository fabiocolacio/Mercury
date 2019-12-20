defmodule Mercury.Repo.Migrations.CreateTables do
  use Ecto.Migration

  def change do
    create table(:users) do
      add :login_key, :binary
      add :challenge, :binary
    end

    create table(:rooms) do

    end

    create table(:participants) do
      add :user_id, references(:users)
      add :room_id, references(:rooms)
    end

    create unique_index(:participants, [:user_id, :room_id])
    
    create table(:messages) do
      add :user_id, references(:users)
      add :room_id, references(:rooms)
      add :timestamp, :utc_datetime
      add :contents, :binary
    end
  end
end
