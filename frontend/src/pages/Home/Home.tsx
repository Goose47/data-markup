import { block } from "../../utils/block";
import "./Home.scss";

const b = block("home");

export const Home = () => {
  return (
    <div className={b()}>
      <h1>Тут может быть общая статистика по ассессорам</h1>
      <p>Но пока здесь ничего нет 🐈</p>
    </div>
  );
};
