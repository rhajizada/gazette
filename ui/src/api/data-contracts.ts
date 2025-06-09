/* eslint-disable */
/* tslint:disable */
// @ts-nocheck
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

export interface GithubComRhajizadaGazetteInternalServiceAddItemToCollectionResponse {
  added_at?: string;
}

export interface GithubComRhajizadaGazetteInternalServiceCollection {
  created_at?: string;
  id?: string;
  last_updated?: string;
  name?: string;
}

export interface GithubComRhajizadaGazetteInternalServiceFeed {
  authors?: GithubComRhajizadaGazetteInternalServicePerson[];
  categories?: string[];
  copyright?: string;
  created_at?: string;
  description?: string;
  feed_link?: string;
  feed_type?: string;
  feed_version?: string;
  generator?: string;
  id?: string;
  image?: any;
  language?: string;
  last_updated_at?: string;
  link?: string;
  links?: string[];
  published_parsed?: string;
  subscribed?: boolean;
  subscribed_at?: string;
  title?: string;
  updated_parsed?: string;
}

export interface GithubComRhajizadaGazetteInternalServiceItem {
  authors?: GithubComRhajizadaGazetteInternalServicePerson[];
  categories?: string[];
  content?: string;
  created_at?: string;
  description?: string;
  enclosures?: any;
  feed_id?: string;
  guid?: string;
  id?: string;
  image?: any;
  liked?: boolean;
  liked_at?: string;
  link?: string;
  links?: string[];
  published_parsed?: string;
  title?: string;
  updated_at?: string;
  updated_parsed?: string;
}

export interface GithubComRhajizadaGazetteInternalServiceLikeItemResponse {
  liked_at?: string;
}

export interface GithubComRhajizadaGazetteInternalServiceListCategoriesResponse {
  categories?: string[];
  limit?: number;
  offset?: number;
  total_count?: number;
}

export interface GithubComRhajizadaGazetteInternalServiceListCollectionsResponse {
  collections?: GithubComRhajizadaGazetteInternalServiceCollection[];
  limit?: number;
  offset?: number;
  total_count?: number;
}

export interface GithubComRhajizadaGazetteInternalServiceListFeedsResponse {
  feeds?: GithubComRhajizadaGazetteInternalServiceFeed[];
  limit?: number;
  offset?: number;
  total_count?: number;
}

export interface GithubComRhajizadaGazetteInternalServiceListItemsResponse {
  items?: GithubComRhajizadaGazetteInternalServiceItem[];
  limit?: number;
  offset?: number;
  total_count?: number;
}

export interface GithubComRhajizadaGazetteInternalServicePerson {
  /** example: jane@example.com */
  email?: string;
  /** example: Jane Doe */
  name?: string;
}

export interface GithubComRhajizadaGazetteInternalServiceSubscibeToFeedResponse {
  subscribed_at?: string;
}

export interface GithubComRhajizadaGazetteInternalServiceUser {
  createdAt?: string;
  email?: string;
  id?: string;
  lastUpdatedAt?: string;
  name?: string;
  sub?: string;
}

export interface InternalHandlerCreateCollectionRequest {
  name?: string;
}

export interface InternalHandlerCreateFeedRequest {
  feed_url?: string;
}
