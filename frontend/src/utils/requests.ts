import axios from "axios";
import { CreateMarkupTypeRq } from "./types";

const API_PREFIX = "https://api.rwfshr.ru";

export const handleCreateMarkupType = async (request: CreateMarkupTypeRq) => {
  // todo: catch
  return await axios
    .post(API_PREFIX + "/api/v1/markupTypes", request)
    .then((data) => {
      return data.data;
    });
};
