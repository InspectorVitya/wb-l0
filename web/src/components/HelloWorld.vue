<template>
  <v-container>
    <v-row >


      <v-col class="mb-4">
        <v-text-field
        v-model="orderId"
        label="Order id"
        hide-details="auto"
      ></v-text-field>
      <v-btn
        class="mt-5"
        block
        @click="getOrder"
        elevation="2"
        large
      >
      Search
      </v-btn>
      </v-col>

      <v-col
        class="mb-5"
        cols="12"
      >
      <div>
        <json-viewer :value="jsonData" :expand-depth=5
  copyable
  boxed></json-viewer>
      </div>
      </v-col>


    </v-row>
  </v-container>
</template>

<script>
  export default {
    name: 'HelloWorld',
    components: {
  },
    data: () => ({
      orderId: "",
      jsonData: null
    }),
    methods: {
      getOrder(){
        console.log(this.orderId);
        this.$http.get("/order/" + this.orderId)
        .then(res => {
          this.jsonData = res.data
        })
        .catch(e =>{
          this.jsonData = e.response.data
        })

      }
    }
  }
</script>

