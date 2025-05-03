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

import {
  GithubComRhajizadaGazetteInternalServiceAddItemToCollectionResponse,
  GithubComRhajizadaGazetteInternalServiceCollection,
  GithubComRhajizadaGazetteInternalServiceFeed,
  GithubComRhajizadaGazetteInternalServiceItem,
  GithubComRhajizadaGazetteInternalServiceLikeItemResponse,
  GithubComRhajizadaGazetteInternalServiceListCollectionItemsResponse,
  GithubComRhajizadaGazetteInternalServiceListCollectionsResponse,
  GithubComRhajizadaGazetteInternalServiceListFeedsResponse,
  GithubComRhajizadaGazetteInternalServiceListItemsResponse,
  GithubComRhajizadaGazetteInternalServiceSubscibeToFeedResponse,
  GithubComRhajizadaGazetteInternalServiceUser,
  InternalHandlerCreateCollectionRequest,
  InternalHandlerCreateFeedRequest,
} from "./data-contracts";
import { ContentType, HttpClient, RequestParams } from "./http-client";

export class Api<
  SecurityDataType = unknown,
> extends HttpClient<SecurityDataType> {
  /**
   * @description Retrieves paginated collections for the current user.
   *
   * @tags Collections
   * @name CollectionsList
   * @summary List collections
   * @request GET:/api/collections
   * @secure
   */
  collectionsList = (
    query: {
      /** Max number of collections */
      limit: number;
      /** Number of collections to skip */
      offset: number;
    },
    params: RequestParams = {},
  ) =>
    this.request<
      GithubComRhajizadaGazetteInternalServiceListCollectionsResponse,
      string
    >({
      path: `/api/collections`,
      method: "GET",
      query: query,
      secure: true,
      ...params,
    });
  /**
   * @description Creates a named collection for the current user.FeedURL@Tags         Collections
   *
   * @name CollectionsCreate
   * @summary Create collection
   * @request POST:/api/collections
   * @secure
   */
  collectionsCreate = (
    body: InternalHandlerCreateCollectionRequest,
    params: RequestParams = {},
  ) =>
    this.request<GithubComRhajizadaGazetteInternalServiceCollection, string>({
      path: `/api/collections`,
      method: "POST",
      body: body,
      secure: true,
      type: ContentType.Json,
      ...params,
    });
  /**
   * @description Retrieves a collection by ID.
   *
   * @tags Collections
   * @name CollectionsDetail
   * @summary Get collection
   * @request GET:/api/collections/{collectionID}
   * @secure
   */
  collectionsDetail = (collectionId: string, params: RequestParams = {}) =>
    this.request<GithubComRhajizadaGazetteInternalServiceCollection, string>({
      path: `/api/collections/${collectionId}`,
      method: "GET",
      secure: true,
      ...params,
    });
  /**
   * @description Deletes a collection by ID.
   *
   * @tags Collections
   * @name CollectionsDelete
   * @summary Delete collection
   * @request DELETE:/api/collections/{collectionID}
   * @secure
   */
  collectionsDelete = (collectionId: string, params: RequestParams = {}) =>
    this.request<void, string>({
      path: `/api/collections/${collectionId}`,
      method: "DELETE",
      secure: true,
      ...params,
    });
  /**
   * @description Adds an item to the specified collection.
   *
   * @tags Collections
   * @name CollectionsItemCreate
   * @summary Add item to collection
   * @request POST:/api/collections/{collectionID}/item/{itemID}
   * @secure
   */
  collectionsItemCreate = (
    collectionId: string,
    itemId: string,
    params: RequestParams = {},
  ) =>
    this.request<
      GithubComRhajizadaGazetteInternalServiceAddItemToCollectionResponse,
      string
    >({
      path: `/api/collections/${collectionId}/item/${itemId}`,
      method: "POST",
      secure: true,
      ...params,
    });
  /**
   * @description Removes the specified item from the collection.
   *
   * @tags Collections
   * @name CollectionsItemDelete
   * @summary Remove item from collection
   * @request DELETE:/api/collections/{collectionID}/item/{itemID}
   * @secure
   */
  collectionsItemDelete = (
    collectionId: string,
    itemId: string,
    params: RequestParams = {},
  ) =>
    this.request<void, string>({
      path: `/api/collections/${collectionId}/item/${itemId}`,
      method: "DELETE",
      secure: true,
      ...params,
    });
  /**
   * @description Retrieves items in the collection, including like status.
   *
   * @tags Collections
   * @name CollectionsItemsList
   * @summary List items in collection
   * @request GET:/api/collections/{collectionID}/items
   * @secure
   */
  collectionsItemsList = (
    collectionId: string,
    query: {
      /** Max number of items */
      limit: number;
      /** Number of items to skip */
      offset: number;
    },
    params: RequestParams = {},
  ) =>
    this.request<
      GithubComRhajizadaGazetteInternalServiceListCollectionItemsResponse,
      string
    >({
      path: `/api/collections/${collectionId}/items`,
      method: "GET",
      query: query,
      secure: true,
      ...params,
    });
  /**
   * @description Returns all feeds or only those the user is subscribed to.
   *
   * @tags Feeds
   * @name FeedsList
   * @summary List feeds
   * @request GET:/api/feeds
   * @secure
   */
  feedsList = (
    query: {
      /** Only subscribed feeds */
      subscribedOnly?: boolean;
      /** Max number of feeds to return */
      limit: number;
      /** Number of feeds to skip */
      offset: number;
    },
    params: RequestParams = {},
  ) =>
    this.request<
      GithubComRhajizadaGazetteInternalServiceListFeedsResponse,
      string
    >({
      path: `/api/feeds`,
      method: "GET",
      query: query,
      secure: true,
      ...params,
    });
  /**
   * @description Creates a new feed from URL or subscribes the user to it.
   *
   * @tags Feeds
   * @name FeedsCreate
   * @summary Create or subscribe feed
   * @request POST:/api/feeds
   * @secure
   */
  feedsCreate = (
    body: InternalHandlerCreateFeedRequest,
    params: RequestParams = {},
  ) =>
    this.request<GithubComRhajizadaGazetteInternalServiceFeed, string>({
      path: `/api/feeds`,
      method: "POST",
      body: body,
      secure: true,
      type: ContentType.Json,
      ...params,
    });
  /**
   * @description Retrieves a feed by ID, including the user’s subscription status.
   *
   * @tags Feeds
   * @name FeedsDetail
   * @summary Get feed
   * @request GET:/api/feeds/{feedID}
   * @secure
   */
  feedsDetail = (feedId: string, params: RequestParams = {}) =>
    this.request<GithubComRhajizadaGazetteInternalServiceFeed, string>({
      path: `/api/feeds/${feedId}`,
      method: "GET",
      secure: true,
      ...params,
    });
  /**
   * @description Removes a feed and all its subscriptions/items.
   *
   * @tags Feeds
   * @name FeedsDelete
   * @summary Delete feed
   * @request DELETE:/api/feeds/{feedID}
   * @secure
   */
  feedsDelete = (feedId: string, params: RequestParams = {}) =>
    this.request<void, string>({
      path: `/api/feeds/${feedId}`,
      method: "DELETE",
      secure: true,
      ...params,
    });
  /**
   * @description Retrieves feed items.
   *
   * @tags Items
   * @name FeedsItemsList
   * @summary List feed items
   * @request GET:/api/feeds/{feedID}/items
   * @secure
   */
  feedsItemsList = (
    feedId: string,
    query: {
      /** Max number of items */
      limit: number;
      /** Number of items to skip */
      offset: number;
    },
    params: RequestParams = {},
  ) =>
    this.request<
      GithubComRhajizadaGazetteInternalServiceListItemsResponse,
      string
    >({
      path: `/api/feeds/${feedId}/items`,
      method: "GET",
      query: query,
      secure: true,
      ...params,
    });
  /**
   * @description Subscribes the user to an existing feed.
   *
   * @tags Feeds
   * @name FeedsSubscribeUpdate
   * @summary Subscribe to feed
   * @request PUT:/api/feeds/{feedID}/subscribe
   * @secure
   */
  feedsSubscribeUpdate = (feedId: string, params: RequestParams = {}) =>
    this.request<
      GithubComRhajizadaGazetteInternalServiceSubscibeToFeedResponse,
      string
    >({
      path: `/api/feeds/${feedId}/subscribe`,
      method: "PUT",
      secure: true,
      ...params,
    });
  /**
   * @description Removes the user’s subscription to a feed.
   *
   * @tags Feeds
   * @name FeedsSubscribeDelete
   * @summary Unsubscribe from feed
   * @request DELETE:/api/feeds/{feedID}/subscribe
   * @secure
   */
  feedsSubscribeDelete = (feedId: string, params: RequestParams = {}) =>
    this.request<void, string>({
      path: `/api/feeds/${feedId}/subscribe`,
      method: "DELETE",
      secure: true,
      ...params,
    });
  /**
   * @description Retrieves items liked by the user, paginated.
   *
   * @tags Items
   * @name ItemsList
   * @summary List liked items
   * @request GET:/api/items
   * @secure
   */
  itemsList = (
    query: {
      /** Max number of items */
      limit: number;
      /** Number of items to skip */
      offset: number;
    },
    params: RequestParams = {},
  ) =>
    this.request<
      GithubComRhajizadaGazetteInternalServiceListItemsResponse,
      string
    >({
      path: `/api/items`,
      method: "GET",
      query: query,
      secure: true,
      ...params,
    });
  /**
   * @description Retrieves an item by ID, including like status for the user.
   *
   * @tags Items
   * @name ItemsDetail
   * @summary Get item
   * @request GET:/api/items/{itemID}
   * @secure
   */
  itemsDetail = (itemId: string, params: RequestParams = {}) =>
    this.request<GithubComRhajizadaGazetteInternalServiceItem, string>({
      path: `/api/items/${itemId}`,
      method: "GET",
      secure: true,
      ...params,
    });
  /**
   * @description Retrieves list of collections that given item is in.
   *
   * @tags Items
   * @name ItemsCollectionsList
   * @summary Get collections item is in.
   * @request GET:/api/items/{itemID}/collections
   * @secure
   */
  itemsCollectionsList = (
    itemId: string,
    query: {
      /** Max number of items */
      limit: number;
      /** Number of items to skip */
      offset: number;
    },
    params: RequestParams = {},
  ) =>
    this.request<
      GithubComRhajizadaGazetteInternalServiceListCollectionsResponse,
      string
    >({
      path: `/api/items/${itemId}/collections`,
      method: "GET",
      query: query,
      secure: true,
      ...params,
    });
  /**
   * @description Creates a like record for the current user on an item.
   *
   * @tags Items
   * @name ItemsLikeCreate
   * @summary Like item
   * @request POST:/api/items/{itemID}/like
   * @secure
   */
  itemsLikeCreate = (itemId: string, params: RequestParams = {}) =>
    this.request<
      GithubComRhajizadaGazetteInternalServiceLikeItemResponse,
      string
    >({
      path: `/api/items/${itemId}/like`,
      method: "POST",
      secure: true,
      ...params,
    });
  /**
   * @description Deletes the like record for the current user on an item.
   *
   * @tags Items
   * @name ItemsLikeDelete
   * @summary Unlike item
   * @request DELETE:/api/items/{itemID}/like
   * @secure
   */
  itemsLikeDelete = (itemId: string, params: RequestParams = {}) =>
    this.request<void, string>({
      path: `/api/items/${itemId}/like`,
      method: "DELETE",
      secure: true,
      ...params,
    });
  /**
   * @description Retrieves currently logged in user.
   *
   * @tags Users
   * @name UserList
   * @summary Get user
   * @request GET:/api/user
   * @secure
   */
  userList = (params: RequestParams = {}) =>
    this.request<GithubComRhajizadaGazetteInternalServiceUser, string>({
      path: `/api/user`,
      method: "GET",
      secure: true,
      ...params,
    });
}
