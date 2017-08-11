defmodule ThreeSum do

  def find_answers(list) do
    {:ok, agent} = Agent.start_link fn -> MapSet.new end
    count = list |> Enum.count
    for i <- 0..count do
      for j <-(i + 1)..count do
        for k <- (j + 1)..count do
          candidate = Enum.map([i,j,k], &(Enum.at(list, &1))) |> Enum.filter(&(&1))
          if Enum.count(candidate) == 3 && Enum.reduce(candidate, 0, &(&1 + &2)) == 0 do
            Agent.update(agent, fn(set) ->
              MapSet.put set, candidate
            end)
          end
        end
      end
    end
    answers = Agent.get(agent, &(&1))
    Agent.stop(agent, :normal)
    answers
  end

  def split_ints(line) do
    line
    |> String.trim_trailing
    |> String.split(" ", trim: true)
    |> Enum.map(&String.to_integer/1)
    |> Enum.sort
  end

  def printf(answers) do
    Enum.each(answers, fn(answer)-> IO.puts Enum.join(answer, " ") end)
    IO.puts("")
  end

  def process_line do
    case IO.read(:line) do
      :eof -> nil
      line -> line
              |> split_ints
              |> find_answers
              |> printf
              process_line()
    end
  end
end

ThreeSum.process_line