import axios from "axios";
import { MarkupTypeRq } from "./types";

const API_PREFIX = "https://api.rwfshr.ru";

export const handleCreateMarkupType = async (request: MarkupTypeRq) => {
  // todo: catch
  return await axios
    .post(API_PREFIX + "/api/v1/markupTypes", request)
    .then((response) => {
      return response.data;
    });
};

export const handleEditMarkupType = async (
  markupId: string,
  request: MarkupTypeRq
) => {
  // todo: catch
  return await axios
    .put(API_PREFIX + "/api/v1/markupTypes/" + markupId, request)
    .then((response) => {
      return response.data;
    });
};

export const getAvailableMarkupTypes = async () => {
  return await axios
    .get(API_PREFIX + "/api/v1/markupTypes")
    .then((response) => response.data.data);
};

export const getDetailedMarkupType = async (markupId: number) => {
  return await axios
    .get(API_PREFIX + "/api/v1/markupTypes/" + markupId)
    .then((response) => response.data);
};

export const deleteMarkupType = async (markupId: number) => {
  return await axios
    .get(API_PREFIX + "/api/v1/markupTypes/" + markupId)
    .then((response) => response.data);
};
