<template>
  <main>
    <section class="section">
      <h1 class="is-size-1">
        Find a Movie
      </h1>
      <div class="columns">
        <b-field
          class="column is-four-fifths pb-0 mb-0"
          label="What do you want to watch?"
        >
          <b-autocomplete
            :data="data"
            placeholder="e.g. The Princess Bride"
            field="title"
            :loading="isFetching"
            @typing="getAsyncData"
            @select="option => selected = option"
          >
            <template slot-scope="props">
              <div class="media">
                <div class="media-left">
                  <img
                    v-if="props.option.poster_path"
                    width="32"
                    :src="`https://image.tmdb.org/t/p/w500/${props.option.poster_path}`"
                  >
                  <b-icon
                    v-else-if="props.option.media_type==='movie'"
                    icon="filmstrip"
                    size="is-medium"
                  />
                  <b-icon
                    v-else-if="props.option.media_type==='tv'"
                    icon="television"
                    size="is-medium"
                  />
                </div>
                <div class="media-content">
                  <strong>{{ props.option.title || props.option.name }}</strong>
                  <br>
                  <small>
                    <span v-if="getYear(props.option.release_date || props.option.first_air_date)">
                      {{ getYear(props.option.release_date || props.option.first_air_date) }}
                    </span>
                  </small>
                </div>
              </div>
            </template>
          </b-autocomplete>
        </b-field>
        <b-field
          class="column is-one-fifth pb-0 mb-0"
          label="Ignore Notion?"
        >
          <b-switch
            v-model="ignoreNotion"
            type="is-success"
          />
        </b-field>
      </div>

      <div
        v-if="selected"
      >
        <div class="buttons mb-5">
          <b-button
            type="is-success has-text-weight-bold"
            expanded
            @click="saveItem()"
          >
            Save
          </b-button>
        </div>
        <div class="columns mb-0">
          <div class="column is-one-third">
            <span
              :class="{'is-clickable' : selected.posters.length > 1}"
              @click="selected.posters.length > 1 ? openPosterModal() : null"
            >
              <b-image
                v-if="selectedPoster"
                :src="`${selectedPoster}`"
                :alt="`Poster for ${selected.title}`"
              />
              <b-icon
                v-else-if="selected.media_type==='movie'"
                icon="filmstrip"
                size="is-large"
              />
              <b-icon
                v-else-if="selected.media_type==='tv'"
                icon="television"
                size="is-large"
              />
            </span>
          </div>
          <div class="column is-two-thirds">
            <h2 class="is-size-2 movie-title level level-left mb-1">
              <span
                v-if="selectedPoster"
                class="mr-2"
              >
                <b-icon
                  v-if="selected.media_type==='movie'"
                  icon="filmstrip"
                  size="is-medium"
                />
                <b-icon
                  v-else-if="selected.media_type==='tv'"
                  icon="television"
                  size="is-medium"
                />
              </span>
              <span>{{ selected.title }}</span>
              <span
                v-if="selected.release_date"
                class="is-size-3"
              >
                &nbsp;({{ selected.release_date }})
              </span>
            </h2>
            <p
              v-if="selected.directors && selected.directors.length"
              class="mt-1"
            >
              <span class="is-size-5">
                {{ makeGrammaticalList(selected.directors) }}
              </span>
            </p>
            <p class="mt-3 mb-4">
              {{ selected.overview }}
            </p>
            <div v-if="(selected.tags || []).length">
              <p class="is-size-4 level mb-1 level">
                Tags
              </p>
              <p>
                <b-tag
                  v-for="(tag, index) in selected.tags"
                  :key="`cast_${index}`"
                  class="m-1"
                  :type="colors[(index+3)%colors.length]"
                >
                  {{ tag }}
                </b-tag>
              </p>
            </div>
            <div v-if="(selected.streamingService || []).length">
              <p class="is-size-4 level mb-1 level">
                Streaming Services
              </p>
              <p>
                <b-tag
                  v-for="(tag, index) in selected.streamingService"
                  :key="`cast_${index}`"
                  class="m-1"
                  :type="colors[(index+4)%colors.length]"
                >
                  {{ tag }}
                </b-tag>
              </p>
            </div>
            <p
              v-if="selected.last_watched"
              class="mt-3 mb-4"
            >
              <strong>
                Last Watched: {{ formatDate(selected.last_watched) }}
              </strong>
            </p>
            <div class="is-align-content-end is-flex is-flex-direction-row-reverse">
              <b-button
                size="is-small"
                class="is-align-content-end level-right is-white"
                @click="openDetailsModal()"
              >
                <span class="has-text-weight-semibold">Edit</span>
              </b-button>
            </div>
          </div>
        </div>
        <div v-if="selected.cast.length">
          <p class="is-size-3 level mb-1 level">
            Cast
            <span class="level-right">
              <span v-if="editingCast">
                <b-button
                  size="is-small"
                  class="is-align-content-end is-white mr-1"
                  @click="editingCast=false"
                >
                  <span class="has-text-weight-semibold">Cancel</span>
                </b-button>
                <b-button
                  size="is-small"
                  class="is-align-content-end is-info"
                  @click="newCast=selectedItem.originalCast"
                >
                  <span class="has-text-weight-semibold">Reset</span>
                </b-button>
                <b-button
                  size="is-small"
                  class="is-align-content-end is-success"
                  @click="saveCast()"
                >
                  <span class="has-text-weight-semibold">Save</span>
                </b-button>

              </span>
              <b-button
                v-else
                size="is-small"
                class="is-align-content-end is-white"
                @click="editCast()"
              >
                <span class="has-text-weight-semibold">Edit</span>
              </b-button>
            </span>
          </p>
          <p v-if="!editingCast">
            <b-tag
              v-for="(person, index) in selected.cast"
              :key="`cast_${index}`"
              class="m-1"
              :type="colors[(index+1)%colors.length]"
            >
              {{ person }}
            </b-tag>
          </p>
          <p v-else>
            <b-field
              type="is-success"
            >
              <b-taginput
                :value="newCast"
                @input="handleTagInput($event, 'newCast')"
              />
            </b-field>
          </p>
        </div>
        <div>
          <p class="is-size-3 level mb-1 mt-2 level">
            Writers
            <span class="level-right">
              <span v-if="editingWriters">
                <b-button
                  size="is-small"
                  class="is-align-content-end is-white mr-1"
                  @click="editingWriters=false"
                >
                  <span class="has-text-weight-semibold">Cancel</span>
                </b-button>
                <b-button
                  size="is-small"
                  class="is-align-content-end is-info"
                  @click="newWriters=selectedItem.originalWriters"
                >
                  <span class="has-text-weight-semibold">Reset</span>
                </b-button>
                <b-button
                  size="is-small"
                  class="is-align-content-end is-success"
                  @click="saveWriters()"
                >
                  <span class="has-text-weight-semibold">Save</span>
                </b-button>

              </span>
              <b-button
                v-else
                size="is-small"
                class="is-align-content-end is-white"
                @click="editWriters()"
              >
                <span class="has-text-weight-semibold">Edit</span>
              </b-button>
            </span>
          </p>
          <p v-if="!editingWriters">
            <span v-if="selected.writers.length">
              <b-tag
                v-for="(person, index) in selected.writers"
                :key="`cast_${index}`"
                class="m-1 tag"
                :class="colors[(index+2)%colors.length]"
              >
                {{ person }}
              </b-tag>
            </span>
            <span v-else>
              <b-tag
                class="m-1 tag"
                :class="colors[(2)%colors.length]"
              >
                None
              </b-tag>
            </span>
          </p>
          <p v-else>
            <b-field
              type="is-success"
            >
              <b-taginput
                :value="newWriters"
                @input="handleTagInput($event, 'newWriters')"
              />
            </b-field>
          </p>
        </div>
      </div>
    </section>

    <b-modal
      v-if="selected && selected.posters.length > 1"
      v-model="isPosterModalActive"
      width="80%"
      scroll="keep"
    >
      <div class="card posters-card">
        <header class="card-header">
          <h3 class="is-size-3 mb-2 card-header-title">
            Posters
          </h3>
        </header>

        <div class="card-content">
          <div class="columns mb-0">
            <div class="column is-one-third relative">
              <div class="fixed top">
                <b-image
                  v-if="newPoster"
                  :src="`${newPoster}`"
                  :alt="`Poster for ${selected.title}`"
                />
                <b-icon
                  v-else
                  icon="television"
                />
              </div>
            </div>
            <div class="column is-two-thirds">
              <div
                v-for="(row, i) in Math.ceil(selected.posters.length/perRow)"
                :key="`poster_row_${i}`"
                class="columns"
              >
                <div
                  v-for="(poster_details, j) in selected.posters.slice(i*perRow, i*perRow+perRow)"
                  :key="`poster_row${i}_item_${j}`"
                  class="column is-one-third selectable"
                  :class="{ 'selected-poster' : newPoster === `https://image.tmdb.org/t/p/original${poster_details.file_path}` }"
                  @click="selectPoster(`https://image.tmdb.org/t/p/original${poster_details.file_path}`)"
                >
                  <b-image
                    :src="`https://image.tmdb.org/t/p/original${poster_details.file_path}`"
                    :alt="`Poster for ${selected.title}`"
                  />
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="p-3">
          <span class="level-right fixed bottom">
            <b-button
              class="is-white mr-2"
              @click="isPosterModalActive=false"
            >
              <span class="has-text-weight-semibold">Cancel</span>
            </b-button>
            <b-button
              class="is-align-content-end is-info mr-2"
              @click="resetPoster()"
            >
              <span class="has-text-weight-semibold">Reset</span>
            </b-button>
            <b-button
              class="is-align-content-end is-success"
              @click="savePoster()"
            >
              <span class="has-text-weight-semibold">Save</span>
            </b-button>

          </span>
        </div>
      </div>
    </b-modal>

    <b-modal
      v-if="selected"
      v-model="isDetailsModalActive"
      width="80%"
      scroll="keep"
    >
      <div class="card">
        <header class="card-header">
          <h3 class="is-size-3 mb-2 card-header-title">
            Details
          </h3>
        </header>

        <div class="card-content">
          <b-field label="Title">
            <b-input
              v-model="newTitle"
              placeholder="Title"
              type="text"
              required
            />
          </b-field>
          <div class="columns">
            <div class="column is-one-third">
              <b-field label="Release Date">
                <b-numberinput
                  v-model="newReleaseDate"
                  :controls="false"
                  min="0"
                  :max="new Date().getFullYear()"
                />
              </b-field>
            </div>
            <div class="column is-two-thirds">
              <b-field label="Media Type">
                <b-radio-button
                  v-model="newMediaType"
                  native-value="tv"
                  expanded
                  type="is-info is-light is-outlined"
                >
                  <b-icon icon="television" />
                  <span>TV Series</span>
                </b-radio-button>

                <b-radio-button
                  v-model="newMediaType"
                  native-value="movie"
                  expanded
                  type="is-info is-light is-outlined"
                >
                  <b-icon icon="filmstrip" />
                  <span>Movie</span>
                </b-radio-button>
              </b-field>
            </div>
          </div>
          <b-field
            label="Director/Creator"
            type="is-success"
          >
            <b-taginput
              :value="newDirectors"
              @input="handleTagInput($event, 'newDirectors')"
            />
          </b-field>
          <b-field
            label="Tags"
            type="is-success"
          >
            <b-taginput
              :value="newTags"
              @input="handleTagInput($event, 'newTags')"
            />
          </b-field>
          <b-field
            label="Streaming Services"
            type="is-success"
          >
            <b-taginput
              :value="newStreamingService"
              @input="handleTagInput($event, 'newStreamingService')"
            />
          </b-field>
          <b-field label="Last Watched">
            <b-datepicker
              v-model="newLastWatched"
              placeholder="Click to select..."
              icon="calendar-today"
              :max-date="new Date()"
              :icon-right="selected ? 'close-circle' : ''"
              icon-right-clickable
              trap-focus
              @icon-right-click="newLastWatched = null"
            />
          </b-field>
          <b-field label="Overview">
            <b-input
              v-model="newOverview"
              type="textarea"
              minlength="1"
              maxlength="1000"
              :placeholder="`Describe the ${newMediaType === 'tv' ? 'TV series' : 'movie'}.`"
            />
          </b-field>
        </div>
        <div class="p-3">
          <span class="level-right">
            <b-button
              class="is-white mr-2"
              @click="isDetailsModalActive=false"
            >
              <span class="has-text-weight-semibold">Cancel</span>
            </b-button>
            <b-button
              class="is-align-content-end is-info mr-2"
              @click="resetDetails()"
            >
              <span class="has-text-weight-semibold">Reset</span>
            </b-button>
            <b-button
              class="is-align-content-end is-success"
              @click="saveDetails()"
            >
              <span class="has-text-weight-semibold">Save</span>
            </b-button>

          </span>
        </div>
      </div>
    </b-modal>
  </main>
</template>

<script>
import debounce from 'lodash/debounce'

export default {
  name: 'IndexPage',
  data () {
    return {
      ignoreNotion: false,
      data: [],
      selectedItem: null,
      isFetching: false,
      perRow: 3,
      selectedPoster: null,
      colors: ['is-dark', 'is-primary', 'is-info', 'is-success', 'is-warning', 'is-danger'],
      editingDetails: false,
      editingCast: false,
      editingWriters: false,
      isPosterModalActive: false,
      isDetailsModalActive: false,

      newTags: [],
      newCast: [],
      newWriters: [],
      newPoster: null,
      newTitle: null,
      newReleaseDate: null,
      newLastWatched: null,
      newOverview: null,
      newDirectors: [],
      newStreamingService: [],
      newMediaType: null
    }
  },
  computed: {
    selected: {
      get () {
        return this.selectedItem
      },
      set (val) {
        if (!val) {
          this.selectedItem = null
          return
        }
        this.isFetching = true
        const api = `${process.env.tmdbURL}/${val.media_type}/${val.id}`
        const params = `api_key=${process.env.tmdbAPIKey}`
        this.$axios.get(`${api}?${params}`)
          .then((details) => {
            const credits = this.$axios.get(`${api}/credits?${params}`)
            const images = this.$axios.get(`${api}/images?${params}`)
            return Promise.all([details, credits, images])
          })
          .then((results) => {
            const details = results[0].data
            const credits = results[1].data
            const images = results[2].data
            const posters = []
            for (const key in images) {
              if (key !== 'logo' && images[key].length) {
                posters.push(...images[key])
              }
            }
            const cast = credits.cast
            const crew = credits.crew
            const primary = []
            const directors = crew.filter(person => person.job.toLowerCase() === 'director')
            details.title = details.title || details.name
            const releaseDate = details.release_date || details.first_air_date
            details.media_type = val.media_type
            details.release_date = this.getYear(releaseDate)
            details.poster_path = `https://image.tmdb.org/t/p/original/${details.poster_path}`
            if (val.media_type === 'tv') {
              primary.push(...details.created_by.map(person => person.name))
            } else {
              primary.push(...directors.map(person => person.name))
            }
            const writers = crew
              .filter((person) => {
                return ['writer', 'story', 'screenplay']
                  .includes(person.job.toLowerCase())
              })
              .map(person => person.name)

            const topCast = cast
              .slice(0, 8)
              .map(person => person.name)

            const selectedItem = {
              ...details,
              cast: topCast,
              directors: primary,
              writers,
              posters: posters
                .filter((poster) => {
                  return poster.iso_639_1 === 'en' || !poster.iso_639_1
                }),

              originalCast: [...topCast],
              originalDirectors: [...primary],
              originalWriters: [...writers],
              originalReleaseDate: details.release_date,
              originalMediaType: details.media_type,
              originalOverview: details.overview,
              originalTitle: details.title
            }
            let notionData = []
            notionData = this.$axios.post('https://viewing.rmccarthy.dev/api/get', {
              title: selectedItem.title || '',
              release_date: selectedItem.release_date || 0,
              media_type: selectedItem.media_type || ''
            })

            return Promise.all([selectedItem, notionData])
          })
          .then((results) => {
            const selectedItem = results[0]
            if (!this.ignoreNotion && ((results[1] || {}).data || []).length) {
              const item = results[1].data[0]

              selectedItem.cast = item.properties.Actors.multi_select.map(actor => actor.name)
              selectedItem.originalCast = selectedItem.cast
              selectedItem.directors = item.properties['Director/Creator'].multi_select.map(director => director.name)
              selectedItem.originalDirectors = selectedItem.directors
              selectedItem.release_date = item.properties['Release Year'].number
              selectedItem.originalReleaseDate = selectedItem.release_date

              selectedItem.tags = item.properties.Tags.multi_select.map(tag => tag.name)
              selectedItem.originalTags = selectedItem.tags

              selectedItem.streamingService = item.properties['Streaming Service'].multi_select.map(ss => ss.name)
              selectedItem.originalStreamingService = selectedItem.streamingService

              const date = (item.properties['Last Watched'].date || {}).start || null
              selectedItem.last_watched = date ? new Date(date) : null
              selectedItem.originalLastWatched = selectedItem.last_watched
              selectedItem.poster_path = ((item.cover || {}).external || {}).url

              selectedItem.overview = ((((item.properties.Summary || {}).rich_text || [])[0] || {}).text || {}).content || ''
              selectedItem.originalOverview = selectedItem.overview
            }

            if (((results[1] || {}).data || []).length) {
              selectedItem.notionPage = results[1].data[0]
            }
            this.selectedItem = selectedItem
            this.selectedPoster = this.selectedItem.poster_path
          })
          .catch((error) => {
            this.selectedItem = null
            throw error
          })
          .finally(() => {
            this.isFetching = false
          })
      }
    },
    formattedItem () {
      // Start with existing notion page if it exists
      const item = {
        page: {
          archived: false
        }
      }

      if (this.selectedItem.notionPage) {
        item.id = this.selectedItem.notionPage.id
      }

      // Add URL for poster if available
      if (this.selectedPoster) {
        item.page.cover = {
          type: 'external',
          external: {
            url: this.selectedPoster
          }
        }
      } else {
        item.cover = null
      }
      const type = this.selectedItem.media_type
      const mediaType = type === 'tv' ? 'TV Series' : type === 'movie' ? 'Movie' : null
      item.page.properties = {
        Actors: {
          type: 'multi_select',
          multi_select: (this.selectedItem.cast || []).map(name => ({ name }))
        },
        Summary: {
          type: 'rich_text',
          rich_text: [
            {
              text: {
                content: this.selectedItem.overview
              }
            }
          ]
        },
        Tags: {
          type: 'multi_select',
          multi_select: (this.selectedItem.tags || []).map(name => ({ name }))
        },
        Title: {
          type: 'title',
          title: [
            {
              text: {
                content: this.selectedItem.title
              }
            }
          ]
        },
        Type: {
          type: 'select',
          select: {
            name: mediaType
          }
        },
        Writers: {
          type: 'multi_select',
          multi_select: (this.selectedItem.writers || []).map(name => ({ name }))
        }
      }
      item.page.properties['Director/Creator'] = {
        type: 'multi_select',
        multi_select: (this.selectedItem.directors || []).map(name => ({ name }))
      }
      item.page.properties['Last Watched'] = {
        type: 'date',
        date: null
      }

      item.page.properties['Release Year'] = {
        type: 'number',
        number: this.selectedItem.release_date
      }
      item.page.properties['Streaming Service'] = {
        type: 'multi_select',
        multi_select: (this.selectedItem.streamingService || []).map(name => ({ name }))
      }
      return item
    }
  },
  methods: {
    async saveItem () {
      const path = ((this.selectedItem || {}).notionPage || {}).id ? 'update' : 'create'
      const resp = await this.$axios.post(`https://viewing.rmccarthy.dev/api/${path}`, this.formattedItem)
      this.selectedItem.notionPage = resp.data
    },
    formatDate (date) {
      function join (t, a, s) {
        function format (m) {
          const f = new Intl.DateTimeFormat('en', m)
          return f.format(t)
        }
        return a.map(format).join(s)
      }

      const a = [{ day: 'numeric' }, { month: 'long' }, { year: 'numeric' }]
      const s = join(date, a, ' ')
      return s
    },
    openDetailsModal () {
      this.newTitle = this.selectedItem.title
      this.newReleaseDate = this.selectedItem.release_date
      this.newOverview = this.selectedItem.overview
      this.newDirectors = this.selectedItem.directors
      this.newMediaType = this.selectedItem.media_type
      this.newTags = this.selectedItem.tags
      this.newLastWatched = this.selectedItem.last_watched
      this.newStreamingService = this.selectedItem.streamingService
      this.isDetailsModalActive = true
    },
    resetDetails () {
      this.selectedItem.title = this.selectedItem.originalTitle
      this.selectedItem.release_date = this.selectedItem.originalReleaseDate
      this.selectedItem.overview = this.selectedItem.originalOverview
      this.selectedItem.directors = this.selectedItem.originalDirectors
      this.selectedItem.media_type = this.selectedItem.originalMediaType
      this.selectedItem.tags = this.selectedItem.originalTags
      this.selectedItem.last_watched = this.selectedItem.originalLastWatched
      this.selectedItem.streamingService = this.selectedItem.originalStreamingService
      this.isDetailsModalActive = false
    },
    async saveDetails () {
      this.selectedItem.title = this.newTitle
      this.selectedItem.release_date = this.newReleaseDate
      this.selectedItem.overview = this.newOverview
      this.selectedItem.directors = this.newDirectors
      this.selectedItem.media_type = this.newMediaType
      this.selectedItem.tags = this.newTags
      this.selectedItem.last_watched = this.newLastWatched
      this.selectedItem.streamingService = this.newStreamingService
      await this.saveItem()
      this.isDetailsModalActive = false
    },
    openPosterModal () {
      this.newPoster = this.selectedPoster
      this.isPosterModalActive = true
    },
    resetPoster () {
      this.newPoster = this.selectedItem.poster_path
      this.selectedPoster = this.selectedItem.poster_path
      this.isPosterModalActive = false
    },
    async savePoster () {
      this.selectedPoster = this.newPoster
      await this.saveItem()
      this.isPosterModalActive = false
    },
    editWriters () {
      this.newWriters = this.selectedItem.writers
      this.editingWriters = true
    },
    async saveWriters () {
      this.selectedItem.writers = this.newWriters
      await this.saveItem()
      this.editingWriters = false
    },
    editCast () {
      this.newCast = this.selectedItem.cast
      this.editingCast = true
    },
    async saveCast () {
      this.selectedItem.cast = this.newCast
      await this.saveItem()
      this.editingCast = false
    },
    handleTagInput (arr, field) {
      this[field] = arr
    },
    makeGrammaticalList (arr) {
      arr = [...arr]
      if (arr.length >= 2) {
        const last = `and ${arr.pop()}`
        arr.push(last)
      }
      let list = ''
      if (arr.length >= 3) {
        list = arr.join(', ')
      } else {
        list = arr.join(' ')
      }
      return list
    },
    selectPoster (path) {
      this.newPoster = path
    },
    getYear (date) {
      return new Date(date).getFullYear()
    },
    reset () {
      this.isFetching = false
      this.editingDetails = false
      this.editingCast = false
      this.editingWriters = false
    },
    getAsyncData: debounce(function (name) {
      this.reset()
      if (!name.length) {
        this.data = []
        return
      }
      this.isFetching = true
      this.$axios.get(`${process.env.tmdbURL}/search/multi?api_key=${process.env.tmdbAPIKey}&query=${name}&include_adult=false`)
        .then(({ data }) => {
          this.data = []
          data.results.forEach((item) => {
            if (item.media_type.toLowerCase() !== 'person') {
              this.data.push(item)
            }
          })
        })
        .catch((error) => {
          this.data = []
          throw error
        })
        .finally(() => {
          this.isFetching = false
        })
    }, 500)
  }
}
</script>

<style scoped>
  .posters-card {
    background-color: hsl(0, 0%, 93%);
  }
  .relative {
    position: relative;
  }
  .fixed {
    position: sticky;
    width: 100%;
  }
  .bottom {
    bottom: 0
  }
  .top {
    top: 0
  }
  .posters-card .card-content {
    max-height: 70vh;
    overflow: scroll;
  }
  .movie-title {
    line-height: 1;
  }
  .selectable:hover {
    outline: none;
    box-shadow: inset 0 0 1px 3px rgba(0, 80, 120, 0.4);
    border-radius: 8px;
    transition: all 0.3s ease;
    cursor: pointer;
  }
  .selected-poster {
    box-shadow: inset 0 0 0 4px rgba(21, 156, 228, 0.4);
    border-radius: 8px;
  }

</style>
